#!/bin/sh
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
PROTOCOL=$(jq '.endpoints[0].protocol' config/conf.json -r)
PROCESSES=$(jq '.processes' config/conf.json -r)
LOGGING=$(jq '.logging' config/conf.json -r)


if [ $PROTOCOL = "http" ]; then
  if [ $LOGGING = true ]; then
    $(gunicorn --chdir restful -w $PROCESSES app:app -b 0.0.0.0:5000 --capture-output --log-level debug);
  else
    $(gunicorn --chdir restful -w $PROCESSES app:app -b 0.0.0.0:5000);
  fi
elif [ $PROTOCOL = "grpc" ]; then
  $(cat config/service.proto > service.proto)
  $(python -m grpc_tools.protoc -I. --python_out=./common --grpc_python_out=./common service.proto);
  $(cd grpc && python pre_app.py)
  2to3 common/ -w -n
  # Uninstall the extra apps
  apt remove -y 2to3 wget
  $(cd grpc && python app.py)
fi
