package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ipInfoMyIPProviderName = "ipinfo.io"
)

type ipInfoResponse struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Postal   string `json:"postal"`
	Region   string `json:"region"`
	Timezone string `json:"timezone"`
	Hostname string `json:"hostname"`
	Org      string `json:"org"`
	Loc      string `json:"loc"`
}

type ipInfoMyIPProvider struct {
}

func newIPInfoMyIPProvider() *ipInfoMyIPProvider {
	return &ipInfoMyIPProvider{}
}

func (i *ipInfoMyIPProvider) name() string {
	return ipInfoMyIPProviderName
}

func (i *ipInfoMyIPProvider) myIP(ctx context.Context) (string, *myIPInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := httpGet(ctx, "https://ipinfo.io/json", "ipinfo.io GET")
	if err != nil {
		return "", nil, err
	}

	respJson := ipInfoResponse{}
	err = json.Unmarshal(resp, &respJson)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse ipinfo.io response as JSON, response: %s\nreason: %w", string(resp), err)
	}

	result := ipInfoToMyIPInfo(&respJson)
	return result.IP, result, nil
}

func ipInfoToMyIPInfo(ipInfo *ipInfoResponse) *myIPInfo {
	return &myIPInfo{
		IP: ipInfo.IP,
		Geo: myIPGeoLocationInfo{
			City:     ipInfo.City,
			State:    ipInfo.Region,
			ZipCode:  ipInfo.Postal,
			Country:  ipInfo.Country,
			Timezone: ipInfo.Timezone,
		},
	}
}
