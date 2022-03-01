# Build and upload Docker images

Build docker image for main application
1. Under [model](/model) directory, run:

``` bash
docker build -t app-demo .
```

After creating the docker images, upload them to each of the clusters 'i' by runnning:
``` bash
kind load docker-image app-demo --name={cluster$i}
```

## Dependecies
1. kind
2. tsung
3. istioctl
4. kubectl
5. go (for installation, configuration and basic testing, follow instructions in e.g. [How to Install GoLang (Go Programming Language) in Linux](HTtps://www.tecmint.com/install-go-in-linux/); make sure go environment variables and path are configured accordingly)

## Environment Preparation
1. Make sure the application-generator folder is located under path ~/go_projects/src/ and initialize module by executing go mod init
2. If needed, install go module dependencies, e.g. cobra and yaml
3. Modify any of the input files under the input directory according to the requirements.
4. Generate and deploy kubernetes manifest files by running 'generator.sh' script. It accepts two arguments, path to input file and value for readiness probe in seconds.
  ```bash
  ./generator.sh {input file} {readiness probe}
  ```
5. Modify the necessary files for request generator
    - Change the initial field of json files under the **tsung** directory according the chain configuration.
    - Change the chain_no field of json files under the **tsung** directory according the chain configuration. For example, for first chain it should be **1**
    - Update the request_task_type of json files under the **tsung** directory for assigning user defined task to each microservice in the chain
    - Change server host ip address in conf.xml file with istio-ingress gateway for first microservice in chain.
    - Change the chain json file under the request section in conf.xml to send request to the desired chain. For example, if first chain is targeted it should be **chain1.json**
6. Change Kubernetes context to the main cluster
```bash
kubectl config use-context cluster1
```
## Running
After configuring environment correctly, you can just use the following command the start request generator.
```bash
tsung -f tsung/conf.xml -k start
```
You can observe the performance metrics for both istio and chain by using the dashboards on grafana interface.
To stop traffic generation use´
```bash
tsung stop
```

For more information see [doc folder](generator/doc) and [masther thesis report](http://www.diva-portal.org/smash/record.jsf?pid=diva2%3A1506576&dswid=8090).
