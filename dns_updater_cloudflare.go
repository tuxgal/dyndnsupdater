package main

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

func getDNSRecord(ctx context.Context, api *cloudflare.API, zoneName string, domainName string) (*cloudflare.DNSRecord, error) {
	zid, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain Zone ID for zone %q, reason: %w", zoneName, err)
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
		return nil, fmt.Errorf(
			"failed to list DNS A records for domain %q in zone %q, reason: %w",
			domainName, zoneName, err)
	}
	log.Debugf(
		"Existing DNS Record(s) for domain %q in zone %q:\n%s",
		domainName, zoneName, prettyPrintJSON(records))

	if len(records) != 1 {
		return nil, fmt.Errorf(
			"Expected %d A record for domain name %q, but obtained %d instead\nRecord(s):\n%s",
			1, domainName, len(records), prettyPrintJSON(records))
	}
	return &records[0], nil
}

func updateCloudflareDNSRecord(
	ctx context.Context, apiToken string, zoneName string, domainName string, externalIP string,
) (bool, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return false, fmt.Errorf("failed to initialize API object, reason: %w", err)
	}

	origRecord, err := getDNSRecord(ctx, api, zoneName, domainName)
	if err != nil {
		return false, err
	}

	if origRecord.Content == externalIP {
		log.Debugf(
			" The A record for domain %q in zone %q is already up to date with the desired External IP %q",
			domainName, zoneName, externalIP)
		return false, nil
	}

	record, err := api.UpdateDNSRecord(
		ctx,
		cloudflare.ZoneIdentifier(origRecord.ZoneID),
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
