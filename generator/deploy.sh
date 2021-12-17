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


#applying manifest files to respective clusters
for d in ./k8s/*; do
	echo "applying deployment manifests to ${d##./k8s/}"
	[[ -d "$d" ]] && kubectl apply --prune -f k8s/${d##./k8s/} -l version=${d##./k8s/} --context ${d##./k8s/}
done
