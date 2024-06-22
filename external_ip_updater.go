package main

import (
	"context"
	"fmt"
)

const (
	myIPFromCloudflareMaxRetries = 3
)

func updateExternalIP(ctx context.Context, token string, zone string, domain string, failOnError bool) error {
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

	ipifyIP, err := myIPFromIPify(ctx)
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

	ipInfo, err := myIPFromIPInfo(ctx)
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

	updated, err := updateCloudflareDNSRecord(ctx, token, zone, domain, ip)
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
