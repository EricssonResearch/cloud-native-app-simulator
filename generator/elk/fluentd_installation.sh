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
#!/bin/bash

if [ $# -eq 0 ]
  then
    echo "No arguments supplied. \n"
    echo "Please provide the <public-ip> of your master node.\n"
else
  # Creating virtual environment
  # Installing prerequisites
  ELASTIC_PORT=30001
  $(python3 -m venv venv)
  source venv/bin/activate
  pip3 install kubernetes
  for d in ../k8s/*; do
      echo "Installing Fluentd on ${d##../k8s/}"
      $(python3 ./update-fluentd.py $1 $ELASTIC_PORT ${d##../k8s/} )
  done
  $(rm -r venv)


fi