#
# Copyright 2023 Ericsson AB
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# This is the base image for the application emulator
# It contains the source code, Go compiler and protobuf compiler
# The generator will compile a unique layered image for the current configuration

FROM golang:1.20

# Install protoc
RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Install grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.19 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Copy relevant parts of the source tree to the new source dir
COPY emulator /usr/src/emulator/emulator
COPY model /usr/src/emulator/model
# Delete placeholder files
RUN rm -Rf /usr/src/emulator/emulator/src/generated

WORKDIR /usr/src/emulator

# Create Go workspace
RUN go work init
RUN go work use ./emulator
RUN go work use ./model

# Download as many modules as possible to be shared between compilations
RUN cd emulator && go mod download -x
