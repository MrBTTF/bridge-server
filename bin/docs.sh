#!/usr/bin/env bash

set -Eeuo pipefail 
set -o xtrace

swag init --pd -g pkg/server/server.go