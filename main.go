// Command dyndnsupdater is a tool to dynamically update the specified DNS record with the machine's external IP.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	myIPFromCloudflareMaxRetries = 3
)

func runOnce(failOnError bool) error {
	ip := ""

	cloudflareIP, err := myIPFromCloudflareWithRetries(myIPFromCloudflareMaxRetries)
	if err != nil {
		if failOnError {
			return err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP obtained using Cloudflare: %q", cloudflareIP)
		ip = cloudflareIP
	}

	ipifyIP, err := myIPFromIPify(context.Background())
	if err != nil {
		if failOnError {
			return err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP obtained using ipify.org: %q", ipifyIP)
		if ip == "" {
			log.Warnf("Using External IP obtained from ipify.org instead of Cloudflare")
			ip = ipifyIP
		}
	}

	ipInfo, err := myIPFromIPInfo(context.Background())
	if err != nil {
		if failOnError {
			return err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP info obtained using ipinfo.io:\n%s", prettyPrintJSON(ipInfo))
		if ip == "" {
			log.Warnf("Using External IP obtained from ipinfo.io instead of Cloudflare")
			ip = ipInfo.IP
		}
	}

	if ip == "" {
		return fmt.Errorf("Unable to obtain External IP from any of the sources, skipping DNS record update ...")
	}

	if cloudflareIP != "" {
		if cloudflareIP != ipInfo.IP {
			log.Warnf(
				"Conflicting External IP information between Cloudflare whoami (%s) and ipinfo.io (%s)",
				cloudflareIP, ipInfo.IP)
			log.Warnf("Using the External IP information from Cloudflare whoami as the trusted source for updating ...")
		}
		if cloudflareIP != ipifyIP {
			log.Warnf(
				"Conflicting External IP information between Cloudflare whoami (%s) and ipify.org (%s)",
				cloudflareIP, ipifyIP)
			log.Warnf("Using the External IP information from Cloudflare whoami as the trusted source for updating ...")
		}
	}

	updated, err := updateCloudflareDNSRecord(
		context.Background(), *cloudflareAPIToken, *cloudflareZoneName, *domainName, ip)
	if err != nil {
		return err
	}
	if updated {
		log.Infof("Updated External IP %q in the A record for domain %q", cloudflareIP, *domainName)
	} else {
		log.Infof("External IP %q in the A record for domain %q is already up to date", cloudflareIP, *domainName)
	}

	return nil
}

func run() int {
	if !validateFlags() {
		return 1
	}

	forever := *daemon
	ranAtLeastOnce := false
	result := 1

	for forever || !ranAtLeastOnce {
		startTime := time.Now()
		nextUpdateTime := startTime.Add(*updateFreq)

		if forever {
			log.Infof("Beginning update ...")
		}

		err := runOnce(false)
		if err != nil {
			log.Errorf("Error querying External IP and updating DNS record, reason: %w", err)
		} else {
			endTime := time.Now()
			if !forever {
				result = 0
			}
			log.Infof("Update took %v since beginning at %v", endTime.Sub(startTime), startTime)
		}

		ranAtLeastOnce = true
		if forever {
			time.Sleep(time.Until(nextUpdateTime))
		}
	}

	return result
}

func main() {
	flag.Parse()
	os.Exit(run())
}
