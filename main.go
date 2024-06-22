// Command dyndnsupdater is a tool to dynamically update the specified DNS record with the machine's external IP.
package main

import (
	"context"
	"flag"
	"os"
	"time"
)

const (
	timestampFormat = "2006-01-02T15:04:05.000Z0700"
)

func checkExporterTerminated(ch chan interface{}) bool {
	timeout := time.NewTimer(100 * time.Millisecond)
	defer timeout.Stop()

	select {
	case <-ch:
		return true
	case <-timeout.C:
	}

	return false
}

func run() int {
	if !validateFlags() {
		return 1
	}

	forever := *daemon
	result := 1
	exporterTerminatedCh := make(chan interface{}, 1)

	if forever {
		go startExporter(exporterTerminatedCh, *listenHost, uint32(*listenPort), *metricsUri)
		// Sleep for 200ms to give sufficient time for the metrics exporter
		// http server to start up.
		time.Sleep(time.Duration(200 * time.Millisecond))
	}

	for {
		nextUpdateTime := time.Now().Add(*updateFreq)

		if forever {
			terminated := checkExporterTerminated(exporterTerminatedCh)
			if terminated {
				log.Errorf("Metrics exporter terminated, exiting ...")
				break
			}
			log.Debugf("Metrics exporter is still alive, proceeding with querying and updating the External IP ...")
			log.Infof("Beginning update ...")
		}

		startTime := time.Now()
		ip, err := updateExternalIP(
			context.Background(), *cloudflareAPIToken, *cloudflareZoneName, *domainName, !forever)
		if err != nil {
			log.Errorf("Error querying External IP and updating DNS record, reason: %w", err)
		} else {
			endTime := time.Now()
			if !forever {
				result = 0
			}
			log.Infof("Update took %v since beginning at %s", endTime.Sub(startTime), startTime.Format(timestampFormat))
			log.Infof("Detailed External IP Info:\n%s", prettyPrintJSON(ip))
		}

		if forever {
			time.Sleep(time.Until(nextUpdateTime))
		} else {
			break
		}
	}

	return result
}

func main() {
	flag.Parse()
	os.Exit(run())
}
