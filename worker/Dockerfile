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

FROM python:3.8.0-alpine

RUN mkdir -p /usr/src/app

WORKDIR /usr/src/app

ADD ./requirements.txt /usr/src/app/requirements.txt

RUN pip install -r requirements.txt && \
    wget -O fortio.tgz https://github.com/fortio/fortio/releases/download/v1.6.8/fortio-linux_x64-1.6.8.tgz && \
    tar -xvf fortio.tgz && \
    mv /usr/src/app/usr/share/fortio /usr/share/ && \
    mv /usr/src/app/usr/bin/fortio /usr/local/bin/ && \
    mkdir /usr/src/app/fortio_result && \
    rm -r /usr/src/app/usr

ENV CONF="/usr/src/app/config/conf.json" 
COPY . /usr/src/app

ENTRYPOINT [ "python" , "main.py"]
