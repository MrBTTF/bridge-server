#!/usr/bin/env bash

set -Eeuo pipefail

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o build/app cmd/main.go