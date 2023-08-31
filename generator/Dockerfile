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

# This is a layered image that builds the current source code on top of hydragen-base and copies the binary to a new image
# The new image should be deployed on all clusters and contains generated code for all endpoints specified in the input file

ARG SRCIMAGE
ARG BASEIMAGE

FROM ${SRCIMAGE} as builder

# Copy everything the application generator created
COPY ./generated/ /usr/src/emulator/emulator/src/generated/

WORKDIR /usr/src/emulator

# Run protobuf compiler
RUN protoc -I/usr/src/emulator --go_out=emulator --go_opt=module=application-emulator --go-grpc_out=emulator --go-grpc_opt=module=application-emulator /usr/src/emulator/emulator/src/generated/service.proto
# Compile the binary without cgo
# This makes it possible to run without any dependencies
RUN CGO_ENABLED=0 go build -mod=readonly -work -ldflags "-s -w" -o app-emulator ./emulator

# Install grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.19 && \
    wget -qO/usr/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /usr/bin/grpc_health_probe

FROM ${BASEIMAGE}

# Copy final binary to /usr/bin/app-emulator
COPY --from=builder /usr/src/emulator/app-emulator /usr/bin/app-emulator
COPY --from=builder /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe

ENV CONF=/usr/src/emulator/config/conf.json

EXPOSE 5000
ENTRYPOINT [ "/usr/bin/app-emulator" ]