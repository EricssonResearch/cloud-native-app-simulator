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
    echo "Please provide the path to the description file.\n"
else
  read -p "Automatically push image to clusters? (Y/N): " push

  if [[ $push == "y" ]] || [[ $push == "Y" ]]; then
    cd k8s
    clusters="$(echo *)"
    cd ..
    contexts="$(kubectl config get-contexts --output=name | tr '\n' ' ')"
    nodes=()

    echo "Trying to discover all nodes that need an updated image..."
    # Try every context with every cluster
    for cl in $clusters; do
      for ctx in $contexts; do
        cmd="kubectl get nodes -o custom-columns=:metadata.name,:spec.taints[].effect --no-headers --cluster $cl --context $ctx"
        output="$($cmd 2>&1 > /dev/null)"
        if [[ $? == 0 ]]; then
          ctxnodes="$($cmd | grep -v 'NoSchedule' | cut -d ' ' -f 1 | tr '\n' ' ')"
          for node in $ctxnodes; do nodes+=("$cl/$node"); done
        fi
      done
    done

    echo "1) Kind (development environment)"
    echo "2) Containerd (requires SSH access)"
    echo "3) Other"
    read runtime 

    if [[ $runtime == "1" ]]; then
      ../community/push-image-to-clusters.sh
    else
      echo "The container image has been saved in generated/hydragen-emulator.tar"
      echo "It should be cached in the namespace k8s.io on the following nodes: ${nodes[@]}"
      echo "The command 'crictl images hydragen-emulator' can be used to verify that the image was loaded correctly"
    fi
  fi

  echo ""

  # Check if logging is required
  LOGGING=$(jq '.settings.logging' $1)
  if [ $LOGGING == true ]; then
    # deploy ELK stack
    $(kubectl apply -f ./elk/namespace.yaml)
    echo "kubectl apply -f elk/"
    echo "Deploying Elasticsearch stack with 3 nodes ..."
    #TODO fix no such file output in mac!
    echo "If you see a message like the following, IGNORE it!"
    echo "\t deploy.sh: line 29: statefulset.apps/es-cluster: No such file or directory\n"
    echo "It might take a while ... (up to 5 minutes, enjoy your coffee! :D) "
    $(kubectl apply -f ./elk/)
    while [ "$( kubectl get sts es-cluster -n logging -o=jsonpath='{.status.readyReplicas}')" != 3 ]; do
       sleep 30
       echo "Waiting for Elasticsearch stack to be ready..."
    done
    KIBANA_NODE=$(kubectl get pods -l=app='kibana' -n logging -o jsonpath='{.items[*].spec.nodeName}')
    NODE_IP=$(kubectl get pods -l=app='kibana' -n logging -o jsonpath='{.items[*].status.hostIP}')
    echo "Kibana is deployed in $KIBANA_NODE"
    echo "To browse kibana, just copy the following address or replace the public ip of the node with the <node-ip>:"
    echo "\t http://<node-ip>:30000"
    echo "\t or"
    echo "\t http://$NODE_IP:30000"
  fi
  # Deploy the microservices to clusters
  for d in ./k8s/*; do
    echo "Applying deployment manifests to ${d##./k8s/}"
    [[ -d "$d" ]] && kubectl apply --prune -f k8s/${d##./k8s/} -l version=${d##./k8s/} --context ${d##./k8s/}
  done
fi
