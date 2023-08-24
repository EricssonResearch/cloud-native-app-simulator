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

has_sudo_password=false
sudo_password=""
has_ssh_password=false
ssh_password=""

while getopts ":s:p:n" option; do
  case "${option}" in
    s)
      has_sudo_password=true
      sudo_password="$OPTARG"
      ;;
    p)
      has_ssh_password=true
      ssh_password="$OPTARG"
      ;;
    n)
      has_sudo_password=true
      sudo_password=""
      ;;
    *)
      echo "Usage: $0 -s <sudo password> -p <ssh password> -n"
      echo "Parameters:"
      echo "  -s: Set sudo password to argument"
      echo "  -p: Set ssh password to argument (if sshpass is installed)"
      echo "  -n: Skip sudo password prompt"
      exit 0
      ;;
    esac
done

name="$(hostname -f)/hydragen-emulator"
image="$(docker images $name --format '{{.Repository}}:{{.Tag}}')"

cd "$(git rev-parse --show-toplevel)/generator/k8s"
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

if [[ $has_sudo_password == false ]]; then
  read -s -p "Sudo password (leave blank if '$(whoami)' has administrative access to containerd): " sudo_password
  if [[ -z "$sudo_password" ]]; then
    echo -n "(not using sudo)"
  fi
  echo ""
fi

for node in "${nodes[@]}"; do
  IFS="/" read -r ctx cl name <<< $node
  # https://kubernetes.io/docs/reference/kubectl/cheatsheet/
  jsonpath="{.status.addresses[?(@.type=='InternalIP')].address}"
  ip="$(kubectl get nodes $name --cluster $cl --context $ctx -o jsonpath=$jsonpath)"
  file="/tmp/containerd-import-image.sh"

  # Start ssh in background
  if [[ $has_ssh_password == true ]]; then
    sshpass -p "$ssh_password" ssh -M -S /tmp/containerd-import-ssh-socket -fnNT "$(whoami)@$ip"
  else
    ssh -M -S /tmp/containerd-import-ssh-socket -fnNT "$(whoami)@$ip"
  fi

  # Copy script to remote machine
  scp -o "ControlPath=/tmp/containerd-import-ssh-socket" ../../community/containerd-import-image.sh "$(whoami)@$ip:/tmp/containerd-import-image.sh"
  # Execute script with archive coming from stdin
  ssh -S /tmp/containerd-import-ssh-socket "$(whoami)@$ip" "chmod +x /tmp/containerd-import-image.sh"
  # Add space at the start to prevent password from being saved in bash history
  cat ../generated/hydragen-emulator.tar | ssh -S /tmp/containerd-import-ssh-socket -C "$(whoami)@$ip" " /tmp/containerd-import-image.sh "$sudo_password""
  ssh -S /tmp/containerd-import-ssh-socket "$(whoami)@$ip" "rm /tmp/containerd-import-image.sh"
  # Close ssh session
  ssh -S /tmp/containerd-import-ssh-socket -O exit "$(whoami)@$ip"
done
