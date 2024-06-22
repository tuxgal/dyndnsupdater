package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Returns the JSON formatted string representation of the specified object.
func prettyPrintJSON(x interface{}) string {
	p, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", x)
	}
	return string(p)
}

func httpGet(ctx context.Context, uri string, comment string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating a new HTTP request, reason: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s failed,\nresp: %s\nreason: %w", comment, prettyPrintJSON(resp), err)
	}
	defer resp.Body.Close()

	log.Debugf("%s - Obtained response:\n%s", comment, prettyPrintJSON(resp))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s failed while reading the body,\nresp: %s\nreason: %w", comment, prettyPrintJSON(resp), err)
	}

	respStr := string(body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s failed due to non-success status code: %d\nbody: %s", comment, resp.StatusCode, respStr)
	}

	log.Debugf("%s obtained response: %s", comment, respStr)
	return body, nil
}
