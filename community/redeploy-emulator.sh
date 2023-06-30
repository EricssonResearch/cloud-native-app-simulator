#!/bin/bash
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

# Assume this file is located in a git repository
cd "$(git rev-parse --show-toplevel)"

echo "Rebuilding emulator image"
docker build -t app-demo .
./community/push-image-to-clusters.sh
echo ""

echo "Restarting all deployments in namespace default"
for d in $(kubectl get -n default -o name deployments)
do
	echo "* $d"
	kubectl rollout restart "$d"
done
