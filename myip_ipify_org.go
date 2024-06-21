package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type IPifyResponse struct {
	IP string
}

func myIPFromIPify() (string, error) {
	resp, err := httpGet("https://api64.ipify.org?format=json", 10*time.Second, "ipify.org GET")
	if err != nil {
		return "", err
	}

	result := IPifyResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse ipify.org response as JSON, response: %s\nreason: %w", string(resp[:]), err)
	}

	return result.IP, nil
}
