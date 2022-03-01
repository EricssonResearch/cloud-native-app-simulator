# Build and upload Docker images
Build docker images for main application and worker
Under [model](/model) directory, run:
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
1. Make sure there exist a kubernetes namespace with name edge-namespace with Istio sidecar injection enabled
2. Make sure the application-generator folder is located under path ~/go_projects/src/ and initialize module by executing go mod init
3. If needed, install go module dependencies, e.g. cobra and yaml
4. Deploy InfluxDB by running startup.sh script and providing no of clusters as argument
    ```bash
    ./startup.sh {$no_of_clusters}
    ```
5. Modify any of the service chain files under the chain directory according to the requirements.
6. Modify any of the cluster placement files under the clusters directory according to the requirements.
7. Generate and deploy kubernetes manifest files by running 'generator.sh' script. It accepts three arguments, path to chain file, path to cluster file and value for readiness probe in seconds.
  ```bash
  ./generator.sh {chain file} {cluster file} {readiness probe}
  ```
8. Modify the necessary files for request generator
    - Change the initial field of json files under the **tsung** directory according the chain configuration.
    - Change the chain_no field of json files under the **tsung** directory according the chain configuration. For example, for first chain it should be **1**
    - Update the request_task_type of json files under the **tsung** directory for assigning user defined task to each microservice in the chain
    - Change server host ip address in conf.xml file with istio-ingress gateway for first microservice in chain.
    - Change the chain json file under the request section in conf.xml to send request to the desired chain. For example, if first chain is targeted it should be **chain1.json**
9. Change Kubernetes context to the main cluster
```bash
kubectl config use-context cluster1
```
10. Open the istioctl grafana dashboard, and add custom data source for InfluxDB
  ```bash
  istioctl dashboard grafana
  ```
  - Configure data source for InfluxDB as the following:
    - Name: influxdb
    - URL: http://influxdb.default.svc.cluster.local:8086
    - Database: latency
    - User: root
    - Password: root
  - Add custom dashboard by importing the json file under the **grafana folder**.
## Running
Note: Make sure there exist an Istio gateway and virtual service for the frontend service(s). For an example see under folder ./frontend/
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
