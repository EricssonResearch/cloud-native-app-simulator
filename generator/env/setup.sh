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


kubectl config use-context cluster1

kubectl create namespace istio-system
kubectl create secret generic cacerts -n istio-system \
    --from-file=samples/certs/ca-cert.pem \
    --from-file=samples/certs/ca-key.pem \
    --from-file=samples/certs/root-cert.pem \
    --from-file=samples/certs/cert-chain.pem

kubectl config use-context cluster2

kubectl create namespace istio-system
kubectl create secret generic cacerts -n istio-system \
    --from-file=samples/certs/ca-cert.pem \
    --from-file=samples/certs/ca-key.pem \
    --from-file=samples/certs/root-cert.pem \
    --from-file=samples/certs/cert-chain.pem

kubectl config use-context cluster3

kubectl create namespace istio-system
kubectl create secret generic cacerts -n istio-system \
    --from-file=samples/certs/ca-cert.pem \
    --from-file=samples/certs/ca-key.pem \
    --from-file=samples/certs/root-cert.pem \
    --from-file=samples/certs/cert-chain.pem


#kubectl config use-context cluster4

#kubectl create namespace istio-system
#kubectl create secret generic cacerts -n istio-system \
#    --from-file=samples/certs/ca-cert.pem \
#    --from-file=samples/certs/ca-key.pem \
#    --from-file=samples/certs/root-cert.pem \
#    --from-file=samples/certs/cert-chain.pem



export MAIN_CLUSTER_CTX=cluster1
export REMOTE_CLUSTER_CTX1=cluster2
export REMOTE_CLUSTER_CTX2=cluster3
#export REMOTE_CLUSTER_CTX3=cluster4

export MAIN_CLUSTER_NAME=cluster1
export REMOTE_CLUSTER_NAME1=cluster2
export REMOTE_CLUSTER_NAME2=cluster3
#export REMOTE_CLUSTER_NAME3=cluster4



istioctl --context=${MAIN_CLUSTER_CTX} manifest apply -f istio-main-cluster.yaml
sleep 10
istioctl --context=${REMOTE_CLUSTER_CTX1} manifest apply -f istio-remote1-cluster.yaml
sleep 10
istioctl --context=${REMOTE_CLUSTER_CTX2} manifest apply -f istio-remote2-cluster.yaml
#sleep 10
#istioctl --context=${REMOTE_CLUSTER_CTX3} manifest apply -f istio-remote3-cluster.yaml


#kubectl apply -f cluster-aware-gateway.yaml --context=${MAIN_CLUSTER_CTX}
#kubectl apply -f cluster-aware-gateway.yaml --context=${REMOTE_CLUSTER_CTX1}
#kubectl apply -f cluster-aware-gateway.yaml --context=${REMOTE_CLUSTER_CTX2}
#kubectl apply -f cluster-aware-gateway.yaml --context=${REMOTE_CLUSTER_CTX3}


istioctl x create-remote-secret --name ${REMOTE_CLUSTER_NAME1} --context=${REMOTE_CLUSTER_CTX1} | \
    kubectl apply -f - --context=${MAIN_CLUSTER_CTX}
istioctl x create-remote-secret --name ${REMOTE_CLUSTER_NAME2} --context=${REMOTE_CLUSTER_CTX2} | \
    kubectl apply -f - --context=${MAIN_CLUSTER_CTX}
#istioctl x create-remote-secret --name ${REMOTE_CLUSTER_NAME3} --context=${REMOTE_CLUSTER_CTX3} | \
#    kubectl apply -f - --context=${MAIN_CLUSTER_CTX}


kind load docker-image redis-demo --name=cluster1
kind load docker-image redis-demo --name=cluster2
kind load docker-image redis-demo --name=cluster3
#kind load docker-image redis-demo --name=cluster4


kind load docker-image redis-demo-worker --name=cluster1
kind load docker-image redis-demo-worker --name=cluster2
kind load docker-image redis-demo-worker --name=cluster3
#kind load docker-image redis-demo-worker --name=cluster4
