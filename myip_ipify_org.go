package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ipifyMyIPProviderName = "ipify.org"
)

type ipifyResponse struct {
	IP string `json:"ip"`
}

type ipifyMyIPProvider struct {
}

func newIPifyMyIPProvider() *ipifyMyIPProvider {
	return &ipifyMyIPProvider{}
}

func (i *ipifyMyIPProvider) name() string {
	return ipifyMyIPProviderName
}

func (i *ipifyMyIPProvider) myIP(ctx context.Context) (string, *myIPInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := httpGet(ctx, "https://api64.ipify.org?format=json", "ipify.org GET")
	if err != nil {
		return "", nil, err
	}

	result := ipifyResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse ipify.org response as JSON, response: %s\nreason: %w", string(resp), err)
	}

	return result.IP, nil, nil
}
