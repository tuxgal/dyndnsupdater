package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ipAPIMyIPProviderName = "ipapi.is"
)

type ipAPIASNInfo struct {
	ASN     uint32 `json:"asn"`
	Route   string `json:"route"`
	Descr   string `json:"descr"`
	Country string `json:"country"`
}

type ipAPICompanyInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Domain  string `json:"domain"`
	Network string `json:"network"`
}

type ipAPILocationInfo struct {
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

type ipAPIResponse struct {
	IP       string
	RIR      string
	ASN      ipAPIASNInfo
	Company  ipAPICompanyInfo
	Location ipAPILocationInfo
}

type ipAPIMyIPProvider struct {
}

func newIPAPIMyIPProvider() *ipAPIMyIPProvider {
	return &ipAPIMyIPProvider{}
}

func (i *ipAPIMyIPProvider) name() string {
	return ipAPIMyIPProviderName
}

func (i *ipAPIMyIPProvider) myIP(ctx context.Context) (string, *myIPInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := httpGet(ctx, "https://api.ipapi.is", "ipapi.is GET")
	if err != nil {
		return "", nil, err
	}

	respJson := ipAPIResponse{}
	err = json.Unmarshal(resp, &respJson)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse ipapi.is response as JSON, response: %s\nreason: %w", string(resp), err)
	}
	log.Debugf("IPAPI resp json:\n%s", prettyPrintJSON(respJson))

	result := ipAPIToMyIPInfo(&respJson)
	return result.IP, result, nil
}

func ipAPIToMyIPInfo(ipapi *ipAPIResponse) *myIPInfo {
	return &myIPInfo{
		IP:  ipapi.IP,
		RIR: ipapi.RIR,
		ASN: myIPASNInfo{
			ASN:         ipapi.ASN.ASN,
			Route:       ipapi.ASN.Route,
			Description: ipapi.ASN.Descr,
			Country:     ipapi.ASN.Country,
		},
		Geo: myIPGeoLocationInfo{
			City:        ipapi.Location.City,
			State:       ipapi.Location.State,
			ZipCode:     ipapi.Location.Zip,
			Country:     ipapi.Location.Country,
			CountryCode: ipapi.Location.CountryCode,
			Continent:   ipapi.Location.Continent,
			Latitude:    ipapi.Location.Latitude,
			Longitude:   ipapi.Location.Longitude,
			Timezone:    ipapi.Location.Timezone,
			IsDST:       ipapi.Location.IsDST,
		},
		Provider: myIPNetworkProviderInfo{
			Name:         ipapi.Company.Name,
			ProviderType: ipapi.Company.Type,
			Network:      ipapi.Company.Network,
			Domain:       ipapi.Company.Domain,
		},
	}
}
