package main

import (
	"flag"

	"github.com/tuxdude/zzzlog"
	"github.com/tuxdude/zzzlogi"
)

var (
	log = buildLogger()
)

func buildLogger() zzzlogi.Logger {
	flag.Parse()

	config := zzzlog.NewConsoleLoggerConfig()
	if *debug {
		config.MaxLevel = zzzlog.LvlDebug
	} else {
		config.MaxLevel = zzzlog.LvlInfo
	}
	return zzzlog.NewLogger(config)
}
