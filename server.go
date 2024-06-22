package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	listenAddrFormat  = "%s:%d"
	metricsURLFormat  = "http://%s:%d/%s"
	metricsPathFormat = "/%s"
)

func promHandler() http.Handler {
	// Use a custom registry to get rid of the default set of metrics
	// added by prometheus and have full control.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		// Process information collector.
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		// Go Build Information.
		collectors.NewBuildInfoCollector(),
		// dyndnsupdater metrics collector.
		// dyndnsUpdaterCollector,
		// Disable Go process metrics collector.
		// collectors.NewGoCollector(),
	)
	return promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.HTTPErrorOnError,
		},
	)
}

func startExporter(
	terminatedCh chan interface{},
	listenHost string,
	listenPort uint32,
	path string,
) {
	defer func() { terminatedCh <- nil }()

	listenAddr := fmt.Sprintf(listenAddrFormat, listenHost, listenPort)
	metricsPath := fmt.Sprintf(metricsPathFormat, path)
	metricsURL := fmt.Sprintf(metricsURLFormat, listenHost, listenPort, path)

	// Set up HTTP handler for metrics.
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promHandler())

	// Start listening for HTTP connections.
	server := http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Infof("Starting dyndnsupdater prometheus metrics exporter on %s", metricsURL)
	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("Failed to start dyndnsupdater prometheus metrics exporter: %s", err)
	}
}
