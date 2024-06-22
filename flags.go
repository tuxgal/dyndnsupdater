package main

import (
	"flag"
	"time"
)

var (
	daemon             = flag.Bool("daemon", false, "Runs in daemon mode (i.e. continuous update mode) when set to true")
	updateFreq         = flag.Duration("updateFreq", time.Duration(5)*time.Minute, "How often are the DNS records updated. Relevant only when running in daemon mode. This cannot be lower than 1m")
	debug              = flag.Bool("debug", false, "Log additional debug information")
	dnsApi             = flag.String("dnsApi", "none", "The DNS API to use. If set to 'none', no DNS record updates will be made. Only supported API at the moment is Cloudflare specified using 'cloudflare'")
	domainName         = flag.String("domainName", "", "The domain name whose A record is updated with the dynamically resolved external IP of the current machine")
	cloudflareAPIToken = flag.String("cloudflareApiToken", "", "The Cloudflare scoped API token used for sending the API requests")
	cloudflareZoneName = flag.String("cloudflareZoneName", "", "The Cloudflare DNS Zone name for the domain A record to be updated")
)

// Returns true if a flag was passed in the command line invocation.
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func validateFlags() bool {
	if !isFlagPassed("dnsApi") {
		log.Fatalf("-dnsApi flag must be specified")
		return false
	}
	if *dnsApi != "cloudflare" {
		log.Fatalf("%q is an invalid value for -dnsApi flag. The only supported and valid option at the moment is 'cloudflare'", *dnsApi)
		return false
	}
	if *updateFreq < time.Duration(1)*time.Minute {
		log.Fatalf("-updateFreq must be at least 1m")
		return false
	}
	return true
}
