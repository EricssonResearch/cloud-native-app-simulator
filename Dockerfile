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

WORKDIR /usr/src/app

RUN apt update
RUN apt upgrade -y

COPY . /usr/src/app

RUN go mod download
RUN go build -o /usr/bin/app-emulator ./emulator

ENV CONF=/usr/src/app/config/conf.json

EXPOSE 5000
ENTRYPOINT ["/usr/bin/app-emulator"]