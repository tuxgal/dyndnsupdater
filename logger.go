package main

import (
	"flag"

	"github.com/tuxgal/tuxlog"
	"github.com/tuxgal/tuxlogi"
)

var log = buildLogger()

func buildLogger() tuxlogi.Logger {
	flag.Parse()

	config := tuxlog.NewConsoleLoggerConfig()
	if *debug {
		config.MaxLevel = tuxlog.LvlDebug
	} else {
		config.MaxLevel = tuxlog.LvlInfo
	}
	return tuxlog.NewLogger(config)
}
