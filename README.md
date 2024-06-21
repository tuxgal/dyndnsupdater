# dyndnsupdater

[![Build](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/build.yml/badge.svg)](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/build.yml) [![Tests](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/tests.yml/badge.svg)](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/tests.yml) [![Lint](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/lint.yml/badge.svg)](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/lint.yml) [![CodeQL](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Tuxdude/dyndnsupdater/actions/workflows/codeql-analysis.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/tuxdude/dyndnsupdater)](https://goreportcard.com/report/github.com/tuxdude/dyndnsupdater)

A CLI (written in go) for primarily updating the external IP of the current
machine in the A record for the specified domain name. It can be run in a
loop to perform these updates periodically and export the metrics from this
updater for prometheus to consume. The tool also prints other information
(eg. Geolocation, internet provider, etc.) about the external IP which is
also exposed through the prometheus metrics.

- CH TXT responses from `whoami.cloudflare.` is used as the primary source
  to determine the external IP.
- ipinfo.io and ipify.org are also used as additional sources to validate
  the obtained external IP.
- ipinfo.io provides geolocation and internet provider information.

# Supported DNS APIs for updating the DNS records

Only `Cloudflare` is supported at the moment for DNS management. If other
providers need to be supported, please file an Issue or submit a pull request
with details for further discussion.

# Usage

Just build and run the binary with the necessary flags depending on your
use case.

TODO: Add more details about the flags, invocations, examples, etc.
