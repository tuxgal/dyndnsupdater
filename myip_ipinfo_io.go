package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return nil, fmt.Errorf("ipinfo.io GET failed,\nresp: %s\nreason: %w", prettyPrintJSON(resp), err)
	}
	defer resp.Body.Close()

	log.Debugf("Obtained response:\n%s", prettyPrintJSON(resp))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ipinfo.io GET failed while reading the body,\nresp: %s\nreason: %w", prettyPrintJSON(resp), err)
	}

	jsonStr := string(body[:])
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ipinfo.io GET failed due to non-success status code: %d\nbody: %s", resp.StatusCode, jsonStr)
	}

	log.Debugf("ipinfo.io GET obtained response: %s", jsonStr)

	result := IPInfoResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ipinfo.io response as JSON, response: %s\nreason: %w", jsonStr, err)
	}

	return &result, nil
}
