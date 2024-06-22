package main

import (
	"context"
	"fmt"
)

const (
	myIPFromCloudflareMaxRetries = 3
)

type ExternalIPASNInfo struct {
	ASN         uint32
	Route       string
	Description string
	Country     string
}

type ExternalIPGeoLocationInfo struct {
	City        string
	State       string
	ZipCode     string
	Country     string
	CountryCode string
	Continent   string
	Latitude    float64
	Longitude   float64
	Timezone    string
	IsDST       bool
}

type ExternalIPNetworkProviderInfo struct {
	Name    string
	Type    string
	Network string
	Domain  string
}

type ExternalIPInfo struct {
	IP       string
	RIR      string
	ASN      ExternalIPASNInfo
	Geo      ExternalIPGeoLocationInfo
	Provider ExternalIPNetworkProviderInfo
}

func updateExternalIP(ctx context.Context, token string, zone string, domain string, failOnError bool) (*ExternalIPInfo, error) {
	ip := ""

	cloudflareIP, err := myIPFromCloudflareWithRetries(myIPFromCloudflareMaxRetries)
	if err != nil {
		if failOnError {
			return nil, err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP obtained using Cloudflare: %s", cloudflareIP)
		ip = cloudflareIP
	}

	ipAPIResp, err := myIPFromIPAPI(ctx)
	if err != nil {
		if failOnError {
			return nil, err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP obtained using ipapi.is:   %s", ipAPIResp.IP)
		if ip == "" {
			log.Warnf("Using External IP obtained from ipapi.is instead of Cloudflare")
			ip = ipAPIResp.IP
		}
	}

	ipifyIP, err := myIPFromIPify(ctx)
	if err != nil {
		if failOnError {
			return nil, err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP obtained using ipify.org:  %s", ipifyIP)
		if ip == "" {
			log.Warnf("Using External IP obtained from ipify.org instead of Cloudflare")
			ip = ipifyIP
		}
	}

	ipInfo, err := myIPFromIPInfo(ctx)
	if err != nil {
		if failOnError {
			return nil, err
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("My External IP obtained using ipinfo.io:  %s", ipInfo.IP)
		if ip == "" {
			log.Warnf("Using External IP obtained from ipinfo.io instead of Cloudflare")
			ip = ipInfo.IP
		}
	}

	if ip == "" {
		return nil, fmt.Errorf("Unable to obtain External IP from any of the sources, skipping DNS record update ...")
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
		return nil, err
	}
	if updated {
		log.Infof("Updated External IP %s in the A record for domain %q", ip, domain)
	} else {
		log.Infof("External IP %s in the A record for domain %q is already up to date", ip, domain)
	}

	return toExternalIPInfo(ipAPIResp, ip), nil
}
