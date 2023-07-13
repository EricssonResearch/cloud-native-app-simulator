#!/bin/bash
# Compiles and runs the emulator with the configured service proto
set -e

echo "Downloading any missing modules..."
(cd /usr/src/app/model && go mod download -x)
(cd /usr/src/app/emulator && go mod download -x)

# TODO compile proto
echo "Compiling service $SERVICE_NAME..."
time go build -mod=readonly -work -ldflags "-s -w" -o /tmp/emulator-$SERVICE_NAME /usr/src/app/emulator

echo "Cleaning up after build..."
go clean -cache
go clean -modcache
rm -Rf /tmp/go-build*

echo "Running /tmp/emulator-$SERVICE_NAME..."
/tmp/emulator-$SERVICE_NAME