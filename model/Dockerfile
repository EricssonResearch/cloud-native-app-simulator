#
# Copyright 2021 Ericsson AB
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

FROM python:3.8.0-slim

RUN mkdir -p /usr/src/app
RUN apt update
RUN apt upgrade -y
RUN apt install -y jq \
    wget \
    2to3

RUN apt install -y libbsd-dev \
    libcap-dev \
    libipsec-mb-dev \
    libjudy-dev \
    libkeyutils-dev \
    libsctp-dev \
    libatomic1 \
    zlib1g-dev \
    libkmod-dev \
    libxxhash-dev \
    git \
    build-essential

WORKDIR /usr/src/app

RUN git clone https://github.com/alekodu/stress-ng.git &&\
    cd stress-ng/ &&\
    git checkout cloudsim &&\
    make clean &&\
    make &&\
    mv stress-ng /usr/local/bin/stress-ng &&\
    cd .. &&\
    rm -R stress-ng/

ADD ./requirements.txt /usr/src/app/requirements.txt

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

# download the grpc health probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.8 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

ENV CONF="/usr/src/app/config/conf.json"

COPY . /usr/src/app

EXPOSE 5000
ENTRYPOINT ["/usr/src/app/run.sh"]