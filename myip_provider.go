package main

import "context"

type myIPASNInfo struct {
	ASN         uint32
	Route       string
	Description string
	Country     string
}

type myIPGeoLocationInfo struct {
	City        string
	State       string
	ZipCode     string
	Country     string
	CountryCode string
	Continent   string
	Latitude    float64
	Longitude   float64
	Timezone    string
	IsDST       bool
}

type myIPNetworkProviderInfo struct {
	Name         string
	ProviderType string
	Network      string
	Domain       string
}

type myIPInfo struct {
	IP       string
	RIR      string
	ASN      myIPASNInfo
	Geo      myIPGeoLocationInfo
	Provider myIPNetworkProviderInfo
}

type myIPProvider interface {
	name() string
	myIP(ctx context.Context) (string, *myIPInfo, error)
}
