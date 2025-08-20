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
    while [ "$( kubectl get sts es-cluster -n logging -o=jsonpath='{.status.readyReplicas}')" != 1 ]; do
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
