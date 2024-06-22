package main

import (
	"context"
	"fmt"
)

func getMyExternalIP(ctx context.Context, failOnError bool) (*myIPInfo, error) {
	if len(myIPProviders) == 0 {
		return nil, fmt.Errorf("no My External IP providers configured")
	}

	var result *myIPInfo
	resultIP := ""
	resultSourceOfTruth := ""
	primary := myIPProviders[0].name()

	for i, p := range myIPProviders {
		ip, info, err := p.myIP(ctx)
		if err != nil {
			if failOnError {
				return nil, err
			} else {
				log.Error(err)
			}
		} else {
			log.Infof("My External IP obtained using %s: %s", p.name(), ip)
			if resultIP == "" {
				resultIP = ip
				resultSourceOfTruth = p.name()
				if i != 0 {
					log.Warnf("Using External IP obtained from non-primary provider %s instead of %s", p.name(), primary)
				}
			} else if ip != resultIP {
				log.Warnf(
					"Conflicting External IP information between %s (%s) and %s (%s)",
					resultSourceOfTruth, resultIP, p.name(), ip)
				log.Warnf("Using the External IP information from %s as the trusted source for updating ...",
					resultSourceOfTruth)

			}

			if result == nil && info != nil && (resultIP == "" || resultIP == info.IP) {
				result = info
				log.Debugf("Using detailed External IP info from provider %s", p.name())
			}
		}
	}

	if resultIP == "" {
		return nil, fmt.Errorf("Unable to obtain External IP from any of the sources, skipping DNS record update ...")
	}

	if resultIP != "" && result == nil {
		result = &myIPInfo{IP: resultIP}
		log.Warnf("Could only obtain the IP but no extra information ...")
	}
	return result, nil
}

func queryAndUpdateExternalIP(ctx context.Context, token string, zone string, domain string, failOnError bool) (*myIPInfo, error) {
	ip, err := getMyExternalIP(ctx, failOnError)
	if err != nil {
		return nil, err
	}

	updated, err := updateCloudflareDNSRecord(ctx, token, zone, domain, ip.IP)
	if err != nil {
		return nil, err
	}
	if updated {
		log.Infof("Updated External IP %s in the A record for domain %q", ip.IP, domain)
	} else {
		log.Infof("External IP %s in the A record for domain %q is already up to date", ip.IP, domain)
	}

	return ip, nil
}
