#!/bin/bash
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

name="$(hostname -f)/hydragen-emulator"
image="$(docker images $name --format '{{.Repository}}:{{.Tag}}')"

cd ../generator/k8s
clusters="$(echo *)"
contexts="$(kubectl config get-contexts --output=name | tr '\n' ' ')"

echo "Trying to discover all nodes that need an updated image..."

names=()
nodes=()

# Try every context with every cluster
# TODO: Does not check for the "node" property in configmap
for cl in $clusters; do
  for ctx in $contexts; do
    cmd="kubectl get nodes -o custom-columns=:metadata.name,:spec.taints[].effect --no-headers --cluster $cl --context $ctx"
    output="$($cmd 2>&1)"
    if [[ $? == 0 ]]; then
      ctxnodes="$(echo "$output" | grep -v 'NoSchedule' | cut -d ' ' -f 1 | tr '\n' ' ')"
      for node in $ctxnodes; do
        names+=("$cl/$node")
        nodes+=("$ctx/$cl/$node")
      done
      break 1
    fi
  done
done

echo "Nodes: ${names[@]}"

read -s -p "Sudo password (leave blank if '$(whoami)' has administrative access to containerd): " password
if [[ -z $password ]]; then
  echo -n "(not using sudo)"
fi
echo ""

for node in "${nodes[@]}"; do
  IFS="/" read -r ctx cl name <<< $node
  # https://kubernetes.io/docs/reference/kubectl/cheatsheet/
  jsonpath="{.status.addresses[?(@.type=='InternalIP')].address}"
  ip="$(kubectl get nodes $name --cluster $cl --context $ctx -o jsonpath=$jsonpath)"
  script="$(cat ../../community/containerd-import-image.sh)"
  file="/tmp/containerd-import-image.sh"

  echo "Sending image to $name at $ip..."
  cat ../generated/hydragen-emulator.tar | \
    ssh -C "$ip" \
    "echo "$script" > $file; chmod +x $file; $file "$password"; rm $file"
done
