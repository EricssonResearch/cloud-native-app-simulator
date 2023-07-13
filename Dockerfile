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

FROM golang:1.20

RUN apt update
RUN apt upgrade -y

# Copy relevant parts of the source tree to the new source dir
COPY emulator /usr/src/app/emulator
COPY model /usr/src/app/model

WORKDIR /usr/src/app

# Create Go workspace
RUN go work init
RUN go work use ./emulator
RUN go work use ./model

# Download as many modules as possible to be shared between pods
RUN (cd emulator && go mod download -x)
RUN (cd model && go mod download -x)

# Don't allow any edits to /usr/src/app by Go compiler
RUN chmod -R a-w /usr/src/app

ENV CONF=/usr/src/app/config/conf.json
ENV PROTO=/usr/src/app/config/service.proto

# HTTP at 5000
# gRPC at 5001
EXPOSE 5000 5001
ENTRYPOINT ["/usr/src/app/emulator/run.sh"]