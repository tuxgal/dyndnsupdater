package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type IPInfoResponse struct {
	City     string
	Country  string
	Hostname string
	IP       string
	Loc      string
	Org      string
	Postal   string
	Region   string
	Timezone string
}

func myIPFromIPInfo() (*IPInfoResponse, error) {
	resp, err := httpGet("https://ipinfo.io/json", 10*time.Second, "ipinfo.io GET")
	if err != nil {
		return nil, err
	}

	result := IPInfoResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ipinfo.io response as JSON, response: %s\nreason: %w", string(resp[:]), err)
	}

	return &result, nil
}
