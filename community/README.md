# Cloud-Native-App-Simulator Community

Welcome to the Cloud-Native-App-Simulator community!

This is the starting point for becoming a contributor - improving code, 
improving docs, etc.

- [Introduction](#introduction)
- [Development Environment](#development-enviornment)

## Introduction

Cloud-Native-App-Simulator is an open source microservice benchmark suite 
for generating complete implementations of more generic architectural patterns
for microservices.

## Development Enviornment
This document helps you get started developing code for Cloud-Native-App-Simulator.
If you follow this guide and find a problem, please take a few minutes to update this page.

The Cloud-Native-App-Simulator build system is designed to run with minimal dependencies:
- kind
- docker
- git

These dependencies are need to be set up before building and running the code.
- [Setting Up Docker](#setting-up-docker)
- [Setting Up Kind](#setting-up-kind)


### Setting Up Docker
To use docker to build you will need:
- **docker tools:** To download and install Docker follow [these instructions](https://docs.docker.com/install/).

### Setting Up Kind
To be able to run the *app-demo container* on a sample cluster, we use 
[Kind](https://kind.sigs.k8s.io/docs/user/quick-start/).

- **Installation:** To download and install Kind follow [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/).
- **Setup the clusters:** To setup the clusters, you could simply run the [`kind-setup-clusters.sh`](kind-setup-clusters.sh)
script.
```
#
# This will create multiple Kind clusters (default 2)
# The naming of each cluster is followed by the number (ie, cluster-1, cluster-2, etc.)
# Each of the created clusters has 3 worker nodes and one control plane by default.
#

./kind-setup-clusters.sh [number of clusters (default 3)] [config of each cluster (default kind-cluster-3-nodes.yaml)]
```



