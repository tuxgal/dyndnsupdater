package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type IPifyResponse struct {
	IP string
}

func myIPFromIPify(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := httpGet(ctx, "https://api64.ipify.org?format=json", "ipify.org GET")
	if err != nil {
		return "", err
	}

	result := IPifyResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse ipify.org response as JSON, response: %s\nreason: %w", string(resp), err)
	}

	return result.IP, nil
}
