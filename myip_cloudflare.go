package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"time"

	"github.com/miekg/dns"
)

const (
	cloudflareMyIPProviderName = "cloudflare"
	cloudflareMyIPTarget       = "whoami.cloudflare."
	cloudflareMyIPResolver     = "1.1.1.1:53"
	cloudflareMyIPMaxRetries   = 3
)

type cloudflareMyIPProvider struct {
}

func newCloudflareMyIPProvider() *cloudflareMyIPProvider {
	return &cloudflareMyIPProvider{}
}

func (i *cloudflareMyIPProvider) name() string {
	return cloudflareMyIPProviderName
}

func (i *cloudflareMyIPProvider) myIP(ctx context.Context) (string, *myIPInfo, error) {
	ip, err := myIPFromCloudflareWithRetries(cloudflareMyIPMaxRetries)
	return ip, nil, err
}

func myIPFromCloudflare() (string, bool, error) {
	client := dns.Client{}
	msg := dns.Msg{}

	msg.Id = dns.Id()
	msg.RecursionDesired = false
	msg.Question = []dns.Question{
		{
			Name:   cloudflareMyIPTarget,
			Qtype:  dns.TypeTXT,
			Qclass: dns.ClassCHAOS,
		},
	}

	resp, rtt, err := client.Exchange(&msg, cloudflareMyIPResolver)
	if err != nil {
		timeout := false
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			timeout = opErr.Timeout()
			log.Debugf("CH TXT record query failed due to an OpError, isTimeout: %t\n%s", timeout, prettyPrintJSON(err))
		} else {
			log.Warnf("CH TXT record query failed:\n%s", prettyPrintJSON(err))
		}

		return "", timeout, fmt.Errorf("CH TXT record query failed, reason: %w", err)
	}

	log.Debugf(
		"CH TXT record query for %q to %q took %v. Response:\n%s",
		cloudflareMyIPTarget,
		cloudflareMyIPResolver,
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

	txt, ok := resp.Answer[0].(*dns.TXT)
	if !ok {
		return "", false, fmt.Errorf("Expected type dns.TXT, but rather obtained: %T", resp.Answer[0])
	}

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

		retryAttempt++
		backoff *= 2

		if retryAttempt > maxRetries {
			return ip, err
		}
	}
}
