# Build and upload Docker image for Application

Build docker image for main application
1. Under [model](/model) directory, run:

``` bash
docker build -t app-demo .
```

2. After creating the docker image, upload it to each of the clusters 'i' by runnning:
``` bash
kind load docker-image app-demo --name={cluster$i}
```
# Build and upload Docker image for Tsung Client

Build docker image for tsung client application
1. Under [client](/client) directory, run:

``` bash
docker build -t tsung .
```

2. After creating the docker image, upload it to each of the clusters 'i' by runnning:
``` bash
kind load docker-image tsung --name={cluster$i}
```

## Dependecies
1. kind
2. tsung
3. istioctl
4. kubectl
5. go (for installation, configuration and basic testing, follow instructions in e.g. [How to Install GoLang (Go Programming Language) in Linux](https://www.tecmint.com/install-go-in-linux/); make sure go environment variables and path are configured accordingly)

## Environment Preparation
1. To generate and deploy microservice-based applications, go to the [generator](/generator) directory
2. Make sure the [generator](/generator) directory is located under path ~/go_projects/src/ and initialize module by executing go mod init
3. If needed, install go module dependencies, e.g. cobra and yaml
4. Modify any of the input files under the **input** directory according to your own requirements (see some json examples under the **examples** directory).
5. Generate and deploy kubernetes manifest files by running the 'generator.sh' script. This command can be run under two different modes: (i) 'random' mode which generates a random description file or (ii) 'preset' mode which generates Kubernetes manifest based on a description file in the input directory". Note that this commands generates k8s yaml files which are stored under the **k8s** directory (see some yaml examples under the **examples** directory).
  ```bash
  ./generator.sh {mode} {input file}
  ```

## Using Tsung for Traffic Load Generation
1. Modify the necessary files for request generation with tsung
    - Change the chain json file under the request section in [conf.xml](/client/conf.xml) to send request to the desired frontend service based on the exposed IP. E.g., for Istio-based deployment, update the http url IP based on the ingress gateway IP where the tsung pod is deployed
2. Upon deployment, the pod will automatically start tsung for traffic generation

For more information see [doc folder](generator/doc) and [masther thesis report](http://www.diva-portal.org/smash/record.jsf?pid=diva2%3A1506576&dswid=8090).

## Contribution Guideline
If you are interested in contribution, please visit the [community page](community) to learn about development environment 
setup.

## Need for Logging
To be able to have logging, simply follow the instructions in [Logging](community/Logging.md).