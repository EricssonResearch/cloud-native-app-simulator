# Development Environment

This document helps you get started developing code for HydraGen.
If you follow this guide and find a problem, please take a few minutes to update this page.

HydraGen's build system is designed to run with minimal dependencies:

- kind
- docker
- git

These dependencies need to be set up before building and running the code.

- [Setting Up Docker](#setting-up-docker)
- [Build the worker image](#build-the-worker-image)
- [Setting Up Kind](#setting-up-kind)

## Setting up Docker

To use docker to build required images you will need:

- **docker tools:** To download and install Docker follow [these instructions](https://docs.docker.com/install/).

## Setting up Kind

To be able to run the *hydragen-emulator container* on a sample cluster, we use
[Kind](https://kind.sigs.k8s.io/docs/user/quick-start/).

- **Installation:** To download and install Kind follow [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/).
- **Setup the clusters:** To setup the clusters, you could simply run the [`kind-setup-clusters.sh`](kind-setup-clusters.sh)
script.

```bash
#
# This will create multiple Kind clusters (default 2)
# The naming of each cluster is followed by the number (ie, cluster-1, cluster-2, etc.)
# Each of the created clusters has 3 worker nodes and one control plane by default.
#

cd community
./kind-setup-clusters.sh [number of clusters (default 2)] [config of each cluster (default kind-cluster-3-nodes.yaml)]
```

## Build the base image

The base image, containing code and compilers, needs to be built from the local source code.

```bash
docker build -t "$(hostname -f)/hydragen-base" .
```

By default, HydraGen will use a release image from GitHub Packages as the base when building your image.
To use the development base image instead set this option in the input JSON configuration:

```json
{
    ...
    "settings": {
        "development": true
    },
    ...
}
```

## Pushing the image to a cluster

The generated image needs to be pushed to all clusters after you have run `generator.sh`.

```bash
cd community
./push-image-to-clusters.sh [number of clusters (default 2)]
```

## Logging

To be able to have logging, simply follow the instructions [here](logging.md).
