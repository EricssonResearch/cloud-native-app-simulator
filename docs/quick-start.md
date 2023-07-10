# Quick Start

## Dependencies

1. docker
2. istioctl
3. kubectl
4. go (for installation, configuration and basic testing, follow instructions in e.g. [How to Install GoLang (Go Programming Language) in Linux](https://www.tecmint.com/install-go-in-linux/); make sure go environment variables and path are configured accordingly)
  
## Build and upload Docker images

In the root directory for the repository, run:

``` bash
docker build -t app-demo .
```

After creating the docker image, upload it to each of the clusters _i_ by running:

``` bash
kind load docker-image app-demo --name={cluster$i}
```

## Environment Preparation

### Deploying a Development Environment

Visit this [page](development-environment.md) to learn about development environment setup.

### Generating and Deploying the Application

1. To generate and deploy microservice-based applications, go to the _/generator/_ directory.
2. If needed, install Go module dependencies, e.g. Cobra and yaml. This can be done by running `go mod download`.
3. Modify any of the input files under the _input/_ directory according to your own requirements (see some json examples under the _examples/_ directory).
4. Generate and deploy kubernetes manifest files by running the _generator.sh_ script. This command can be run under two different modes, preset mode or random mode, with the syntax `./generator.sh {mode} {input file}`. See [Application Generator](home.md#application-generator) for more details.
5. Change Kubernetes context to the main cluster: `kubectl config use-context cluster1`.

## Generating Traffic

Modify the necessary files for traffic generation. For example, with Tsung:

- Change the chain json file under the request section in the _conf.xml_ file to send request to the desired frontend service based on the exposed IP, e.g., for Istio-based deployment, update the http url IP based on the ingress gateway IP where the tsung pod is deployed.
- After that, you can just use the following command to start traffic generation: `tsung -f tsung/conf.xml -k start`.
- To stop traffic generation use `tsung stop`.

## Traffic Monitoring

You can observe the performance metrics for the application traffic by using the dashboards on the grafana web UI.
