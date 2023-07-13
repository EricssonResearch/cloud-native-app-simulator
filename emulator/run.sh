#!/bin/bash
# Compiles and runs the emulator with the configured service proto
set -e

echo "Downloading any missing modules..."
(cd /usr/src/app/model && go mod download)
(cd /usr/src/app/emulator && go mod download)

# TODO compile proto
echo "Compiling service $SERVICE_NAME..."
go build -o /tmp/emulator-$SERVICE_NAME /usr/src/app/emulator

echo "Cleaning up after build..."
go clean -cache
go clean -modcache

echo "Running /tmp/emulator-$SERVICE_NAME..."
/tmp/emulator-$SERVICE_NAME