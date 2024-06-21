package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Returns the JSON formatted string representation of the specified object.
func prettyPrintJSON(x interface{}) string {
	p, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", x)
	}
	return string(p)
}

// Returns true if a flag was passed in the command line invocation.
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func httpGet(uri string, timeout time.Duration, comment string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return nil, fmt.Errorf("%s failed,\nresp: %s\nreason: %w", comment, prettyPrintJSON(resp), err)
	}
	defer resp.Body.Close()

	log.Debugf("%s - Obtained response:\n%s", comment, prettyPrintJSON(resp))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s failed while reading the body,\nresp: %s\nreason: %w", comment, prettyPrintJSON(resp), err)
	}

	respStr := string(body[:])
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s failed due to non-success status code: %d\nbody: %s", comment, resp.StatusCode, respStr)
	}

	log.Debugf("%s obtained response: %s", comment, respStr)
	return body, nil
}
