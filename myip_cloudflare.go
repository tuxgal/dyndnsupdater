package main

import (
	"fmt"
	"math/rand/v2"
	"net"
	"time"

	"github.com/miekg/dns"
)

const (
	cloudflareCurrentIPTarget   = "whoami.cloudflare."
	cloudflareCurrentIPResolver = "1.1.1.1:53"
)

func myIPFromCloudflare() (string, bool, error) {
	client := dns.Client{}
	msg := dns.Msg{}

	msg.Id = dns.Id()
	msg.RecursionDesired = false
	msg.Question = []dns.Question{
		{
			Name:   cloudflareCurrentIPTarget,
			Qtype:  dns.TypeTXT,
			Qclass: dns.ClassCHAOS,
		},
	}

	resp, rtt, err := client.Exchange(&msg, cloudflareCurrentIPResolver)
	if err != nil {
		opErr, ok := err.(*net.OpError)
		if ok {
			log.Debugf("CH TXT record query failed due to an OpError, isTimeout: %t\n%s", opErr.Timeout(), prettyPrintJSON(err))
		} else {
			log.Warnf("CH TXT record query failed:\n%s", prettyPrintJSON(err))
		}

		return "", (ok && opErr.Timeout()), fmt.Errorf("CH TXT record query failed, reason: %w", err)
	}

	log.Debugf(
		"CH TXT record query for %q to %q took %v. Response:\n%s",
		cloudflareCurrentIPTarget,
		cloudflareCurrentIPResolver,
		rtt,
		prettyPrintJSON(resp))
	if len(resp.Answer) == 0 {
		return "", false, fmt.Errorf("no results in CH TXT record query response")
	}
	if len(resp.Answer) > 1 {
		log.Warnf(
			"Found %d entries in answer section when only 1 is expected. Using just the first one instead ...",
			len(resp.Answer))
	}
	txt := resp.Answer[0].(*dns.TXT)
	if len(txt.Txt) > 1 {
		log.Warnf(
			"Found %d TXT records in the answer section when only 1 is expected. Using just the first one instead ...",
			len(txt.Txt))
	}
	return txt.Txt[0], false, nil
}

func myIPFromCloudflareWithRetries(maxRetries uint32) (string, error) {
	retryAttempt := uint32(0)
	backoff := 1
	for {
		ip, canRetry, err := myIPFromCloudflare()
		if err == nil || !canRetry {
			return ip, err
		}

		log.Warnf("Obtained retriable error while looking up my IP: %v", err)
		// Exponential backoff with a random offset (min 0ms, max 999ms).
		duration := (time.Duration(backoff) * time.Second) + (time.Duration(rand.IntN(1000)) * time.Millisecond)
		log.Warnf("Backing off %v before retrying ...", duration)
		time.Sleep(duration)

		retryAttempt += 1
		backoff = backoff * 2

		if retryAttempt > maxRetries {
			return ip, err
		}
	}
}
