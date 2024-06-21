package main

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

func updateCloudflareDNSRecord(
	ctx context.Context, apiToken string, zoneName string, domainName string, externalIP string) (bool, error) {
	api, err := cloudflare.NewWithAPIToken(*cloudflareAPIToken)
	if err != nil {
		return false, fmt.Errorf("failed to initialize API object, reason: %w", err)
	}

	zid, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return false, fmt.Errorf("failed to obtain Zone ID for zone %q, reason: %w", zoneName, err)
	}
	log.Debugf("Zone ID for zone %q: %q", zoneName, zid)

	records, _, err := api.ListDNSRecords(
		ctx,
		cloudflare.ZoneIdentifier(zid),
		cloudflare.ListDNSRecordsParams{
			Type: "A",
			Name: domainName,
		})
	if err != nil {
		return false, fmt.Errorf(
			"failed to list DNS A records for domain %q in zone %q, reason: %w",
			domainName, zoneName, err)
	}
	log.Debugf(
		"Existing DNS Record for domain %q in zone %q:\n%s",
		domainName, zoneName, prettyPrintJSON(records))

	if len(records) != 1 {
		return false, fmt.Errorf(
			"Expected %d A records for domain name %q, but obtained %d instead\nRecords:\n%s",
			1, domainName, len(records), prettyPrintJSON(records))
	}
	origRecord := records[0]

	if origRecord.Content == externalIP {
		log.Debugf(
			"A record for domain %q in zone %q is already up to date with the desired External IP %q",
			domainName, zoneName, externalIP)
		return false, nil
	}

	record, err := api.UpdateDNSRecord(
		ctx,
		cloudflare.ZoneIdentifier(zid),
		cloudflare.UpdateDNSRecordParams{
			Type:     "A",
			Name:     origRecord.Name,
			Content:  externalIP,
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
			externalIP, domainName, zoneName, err)
	}

	log.Debugf("Updated DNS Record:\n%s", prettyPrintJSON(record))
	return true, nil
}
