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

DEFAULT_NUM=2
DEFAULT_CONFIG="kind-cluster-3-nodes.yaml"
if [ -z "$1" ]; then
	NUM=$DEFAULT_NUM
else
	NUM=$1
fi

if [ -z "$2" ]; then
  CONFIG=$DEFAULT_CONFIG
else
  CONFIG=$2
fi



# Create the kind multi-node clusters based on the given config
for i in $(seq ${NUM}); do
  kind create cluster --name cluster-${i} --config $CONFIG
done