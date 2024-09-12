package main

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

func getDNSRecord(ctx context.Context, api *cloudflare.API, zone string, domain string) (string, *cloudflare.DNSRecord, error) {
	zid, err := api.ZoneIDByName(zone)
	if err != nil {
		return "", nil, fmt.Errorf("failed to obtain Zone ID for zone %q, reason: %w", zone, err)
	}
	log.Debugf("Zone ID for zone %q: %q", zone, zid)

	records, _, err := api.ListDNSRecords(
		ctx,
		cloudflare.ZoneIdentifier(zid),
		cloudflare.ListDNSRecordsParams{
			Type: "A",
			Name: domain,
		})
	if err != nil {
		return "", nil, fmt.Errorf(
			"failed to list DNS A records for domain %q in zone %q, reason: %w",
			domain, zone, err)
	}
	log.Debugf(
		"Existing DNS Record(s) for domain %q in zone %q:\n%s",
		domain, zone, prettyPrintJSON(records))

	if len(records) != 1 {
		return "", nil, fmt.Errorf(
			"Expected %d A record for domain name %q, but obtained %d instead\nRecord(s):\n%s",
			1, domain, len(records), prettyPrintJSON(records))
	}
	return zid, &records[0], nil
}

func updateCloudflareDNSRecord(
	ctx context.Context, token string, zone string, domain string, ip string,
) (bool, error) {
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return false, fmt.Errorf("failed to initialize API object, reason: %w", err)
	}

	zid, origRecord, err := getDNSRecord(ctx, api, zone, domain)
	if err != nil {
		return false, err
	}

	if origRecord.Content == ip {
		log.Debugf(
			" The A record for domain %q in zone %q is already up to date with the desired External IP %q",
			domain, zone, ip)
		return false, nil
	}

	record, err := api.UpdateDNSRecord(
		ctx,
		cloudflare.ZoneIdentifier(zid),
		cloudflare.UpdateDNSRecordParams{
			Type:     "A",
			Name:     origRecord.Name,
			Content:  ip,
			Data:     origRecord.Data,
			Priority: origRecord.Priority,
			ID:       origRecord.ID,
			TTL:      origRecord.TTL,
			Proxied:  origRecord.Proxied,
			Comment:  &origRecord.Comment,
			Tags:     origRecord.Tags,
		},
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to update DNS A record with IP %q for domain %q in zone %q, reason: %w",
			ip, domain, zone, err)
	}

	log.Debugf("Updated DNS Record:\n%s", prettyPrintJSON(record))
	return true, nil
}
