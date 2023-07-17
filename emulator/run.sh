#!/bin/bash
# Compiles and runs the emulator with the configured service proto
set -e

cd /usr/src/app

echo "Downloading any missing modules..."
(cd model && go mod download -x)
(cd emulator && go mod download -x)

echo "Generating gRPC code for $SERVICE_NAME..."
rm -Rf emulator/src/generated/*
ln -s $GRPCIMPL emulator/src/generated/impl.go
protoc --go_out=emulator --go_opt=module=application-emulator --go-grpc_out=emulator --go-grpc_opt=module=application-emulator $GRPCPROTO

echo "Compiling service $SERVICE_NAME..."
time go build -mod=readonly -work -ldflags "-s -w" -o /tmp/emulator-$SERVICE_NAME emulator

echo "Cleaning up after build..."
go clean -cache
go clean -modcache
rm -Rf /tmp/go-build*

echo "Running /tmp/emulator-$SERVICE_NAME..."
/tmp/emulator-$SERVICE_NAME