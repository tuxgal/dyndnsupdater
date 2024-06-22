package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type IPAPIASNInfo struct {
	ASN     uint32 `json:"asn"`
	Route   string `json:"route"`
	Descr   string `json:"descr"`
	Country string `json:"country"`
}

type IPAPICompanyInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Domain  string `json:"domain"`
	Network string `json:"network"`
}

type IPAPILocationInfo struct {
	City        string  `json:"city"`
	State       string  `json:"state"`
	Zip         string  `json:"zip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Continent   string  `json:"continent"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	IsDST       bool    `json:"is_dst"`
}

type IPAPIResponse struct {
	IP       string
	RIR      string
	ASN      IPAPIASNInfo
	Company  IPAPICompanyInfo
	Location IPAPILocationInfo
}

func myIPFromIPAPI(ctx context.Context) (*IPAPIResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := httpGet(ctx, "https://api.ipapi.is", "ipapi.is GET")
	if err != nil {
		return nil, err
	}

	result := IPAPIResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ipapi.is response as JSON, response: %s\nreason: %w", string(resp), err)
	}

	return &result, nil
}
