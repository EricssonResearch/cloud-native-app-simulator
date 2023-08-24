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

name="$(hostname -f)/hydragen-emulator"
image="$(docker images $name --format '{{.Repository}}:{{.Tag}}')"

cd "$(git rev-parse --show-toplevel)/generator/k8s"
contexts="$(echo *)"
#contexts="$(kubectl config get-contexts --output=name | tr '\n' ' ')"

echo "Contexts: $contexts"
echo "Trying to discover all nodes that need an updated image..."
echo ""

nodes=()

# Try every context
# TODO: Does not check for the "node" property in configmap
for ctx in $contexts; do
  echo "Trying to access context $ctx"
  cmd="kubectl get nodes -o custom-columns=:metadata.name,:spec.taints[].effect --no-headers --context $ctx"
  output="$($cmd 2>&1)"
  if [[ $? == 0 ]]; then
    echo "  Kubectl returned nodes: $(echo $output | tr '\n' ' ')"
    ctxnodes="$(echo "$output" | grep -v 'NoSchedule' | cut -d ' ' -f 1 | tr '\n' ' ')"
    echo "  Nodes that can have pods scheduled: $ctxnodes"
    for node in $ctxnodes; do
      nodes+=("$ctx/$node")
    done
  else
    echo "  Command failed (exit status $?): $(echo $output | tr '\n' ' ')"
  fi
  echo "" 
done

echo "Nodes: ${nodes[@]}"

read -s -p "Sudo password (leave blank if '$(whoami)' has administrative access to containerd): " password
if [[ -z $password ]]; then
  echo -n "(not using sudo)"
fi
echo ""

for node in "${nodes[@]}"; do
  IFS="/" read -r ctx name <<< $node
  # https://kubernetes.io/docs/reference/kubectl/cheatsheet/
  jsonpath="{.status.addresses[?(@.type=='InternalIP')].address}"
  ip="$(kubectl get nodes $name --context $ctx -o jsonpath=$jsonpath)"
  file="/tmp/containerd-import-image.sh"

  # Start ssh in background
  ssh -M -S /tmp/containerd-import-ssh-socket -fnNT "$(whoami)@$ip" 
  # Copy script to remote machine
  scp -o "ControlPath=/tmp/containerd-import-ssh-socket" ../../community/containerd-import-image.sh "$(whoami)@$ip:/tmp/containerd-import-image.sh"
  # Execute script with archive coming from stdin
  ssh -S /tmp/containerd-import-ssh-socket "$(whoami)@$ip" "chmod +x /tmp/containerd-import-image.sh"
  # Add space at the start to prevent password from being saved in bash history
  cat ../generated/hydragen-emulator.tar | ssh -S /tmp/containerd-import-ssh-socket -C "$(whoami)@$ip" " /tmp/containerd-import-image.sh "$password""
  ssh -S /tmp/containerd-import-ssh-socket "$(whoami)@$ip" "rm /tmp/containerd-import-image.sh"
  # Close ssh session
  ssh -S /tmp/containerd-import-ssh-socket -O exit "$(whoami)@$ip"
done
