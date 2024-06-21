package main

import "flag"

var (
	debug              = flag.Bool("debug", false, "Log additional debug information")
	dnsApi             = flag.String("dnsApi", "none", "The DNS API to use. If set to 'none', no DNS record updates will be made. Only supported API at the moment is Cloudflare specified using 'cloudflare'")
	domainName         = flag.String("domainName", "", "The domain name whose A record is updated with the dynamically resolved external IP of the current machine")
	cloudflareAPIToken = flag.String("cloudflareApiToken", "", "The Cloudflare scoped API token used for sending the API requests")
	cloudflareZoneName = flag.String("cloudflareZoneName", "", "The Cloudflare DNS Zone name for the domain A record to be updated")
)
