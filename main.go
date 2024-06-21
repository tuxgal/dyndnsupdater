// Command dyndnsupdater is a tool to dynamically update the specified DNS record with the machine's external IP.
package main

import (
	"context"
	"flag"
	"os"
)

func run() int {
	if !isFlagPassed("dnsApi") {
		log.Fatalf("-dnsApi flag must be specified")
		return 1
	}
	if *dnsApi != "cloudflare" {
		log.Fatalf("%q is an invalid value for -dnsApi flag. The only supported and valid option at the moment is 'cloudflare'", *dnsApi)
		return 1
	}

	ip, err := myIPFromCloudflareWithRetries(3)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	log.Infof("My External IP obtained using Cloudflare: %q", ip)

	ipInfo, err := myIPFromIPInfo()
	if err != nil {
		log.Fatal(err)
		return 1
	}
	log.Infof("External IP info obtained from ipinfo.io:\n%s", prettyPrintJSON(ipInfo))

	if ip != ipInfo.IP {
		log.Warnf(
			"Conflicting External IP information between Cloudflare whoami (%s) and ipinfo.io (%s)",
			ip, ipInfo.IP)
		log.Warnf("Using the External IP information from Cloudflare whoami as the trusted source for updating ...")
	}

	updated, err := updateCloudflareDNSRecord(
		context.Background(), *cloudflareAPIToken, *cloudflareZoneName, *domainName, ip)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	if updated {
		log.Infof("Updated External IP %q in the A record for domain %q", ip, *domainName)
	} else {
		log.Infof("External IP %q in the A record for domain %q is already up to date", ip, *domainName)
	}

	return 0
}

func main() {
	flag.Parse()
	os.Exit(run())
}
