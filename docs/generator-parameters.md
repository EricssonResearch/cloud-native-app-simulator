# Input File

HydraGen uses a standard taxonomy with hierarchical structure, which is the input to the Generator module. The input to this module must be given in JSON format following the taxonomy described here.

## Overall Application Taxonomy

At least one cluster and endpoint is required. Other sections are optional and will be set to their default values if omitted.

#### Required attributes

* **name**: The name of the service object in Kubernetes.
* **protocol**. Determines if endpoints should respond to HTTP or gRPC requests.
* **clusters**: An array of clusters that the service will be deployed on.
* **endpoints**: An array of HTTP/gRPC endpoints that the service exposes.

#### Optional attributes

* **logging**: Enables logging using Elasticsearch. See [logging.md](logging.md) for more information.
* **development**: Builds the application emulator from a local source image (`hydragen-base`) instead of the latest release image.
* **base_image**: Specifies the base Docker image for the application emulator. For example, to use Ubuntu 20.04, set this to `ubuntu:20.04`. The default is `busybox` which provides a minimal shell and set of utilities.
* **resources**: Resource allocation requests and limits.
* **processes**: The maximum number of processes the service is allowed to use (`GOMAXPROCS`). If this is set to 0, the Go runtime will choose the number of processes to use. Default: 0
* **readiness_probe**: The initial delay before readiness probe is initiated. Default: 1 second

#### Format

```json
{
  "settings": {
    "logging": <boolean>,
    "development": <boolean>,
    "base_image": "<string>"
  },
  "services": [
    {
      "name": "<string>",
      "protocol": "<string:http|grpc>",
      "clusters": [...],
      "resources": {...},
      "processes": <integer>,
      "readiness_probe": <integer:seconds>,
      "endpoints": [...]
    }
  ],
  ...
}
```

## Describing Workload Placement and Scaling

The user can define the placement of each microservice on a specific node, cluster, or namespace. It is also possible to specify the number of microservice replicas once deployed.

#### Required attributes

* **cluster**: The cluster that the service will be deployed on.

#### Optional attributes

* **node**: Constrain the service to run on a specific node (for example, "cluster1-worker1"). Default: Empty string
* **namespace**: The namespace that the service will be created in. Default: "default"
* **replicas**: The number of replicas to create. Default: 1
* **annotations**: An array of arbitrary metadata to attach to the service. Default: Empty array

#### Format

```json
"clusters": [
    {
      "cluster": "<string>",
      "node": "<string>",
      "namespace": "<string>",
      "replicas": <integer>,
      "annotations": [...]
    },
    ...
]
```

## Describing Resource Allocation

HydraGen also supports the configuration of the requested resources to be allocated to a microservice instance and the maximum resource usage in terms of both CPU and memory.

#### Optional attributes

* **limits/cpu**: The maximum amount of CPU time that the service can use. Default: 1000m
* **limits/memory**: The maximum amount of memory that the service can use. Default: 1024M
* **requests/cpu**: The desired amount of CPU time for the service. Default: 500m
* **requests/memory**: The desired amount of memory for the service. Default: 256M

```json
"resources": {
  "limits": {
    "cpu": "<string:cores>",
    "memory": "<string:bytes>",
  },
  "requests": {
    "cpu": "<string:cores>",
    "memory": "<string:bytes>"
  }
}
```

## Describing Topological Architecture

For each microservice, HydraGen supports a set of configuration parameters that define the topological architecture of an application by describing the dependencies between services. To define the microservice fan-in, different parameters can be used which specify the set of endpoints a component serves. For each endpoint, the user can specify parameters such as a relative fan-out based on a set of calls to subsequent microservice endpoints as well as the execution mode across these calls. These options enable the user to generate complex multi-tier application architectures with different fan-in and/or fan-out characteristics.

#### Required attributes

* **name**: The request path of the endpoint. Can only contain lowercase alphanumeric characters, '.' or '-'.

#### Optional attributes

* **execution_mode**: Determines if the server responding at this endpoint should handle requests sequentially or in parallel, on multiple threads. Default: "sequential"
* **cpu_complexity**: CPU stress parameters.
* **network_complexity**: Network stress parameters.

#### Format

```json
"endpoints": [
  {
    "name": "<string>",
    "execution_mode": "<string:sequential|parallel>",

    "cpu_complexity": {...},
    "network_complexity": {...}
  },
  ...
]
```

## Describing Resource Stressors

HydraGen supports parameters to express the computational complexity or stress a microservice exerts on the different hardware resources. Initially, CPU-bounded or network-bounded tasks are implemented. The complexity of a CPU-bounded task can be described based on the time a busy-wait is executed, while the load on the network I/O can be described by specifying parameters such as the call forwarding mode and the request/response size for each service endpoint call.

Documentation for implementing a new stressor can be found [here](stressors.md).

### CPU Complexity

The CPU stressor will lock threads for exclusive access while it is executing. This prevents the service from responding to requests on that thread.

#### Required attributes

* **execution_time**: Determines how much time each thread will spend busy-waiting when responding to a request.

#### Optional attributes

* **threads**: The number of threads (goroutines) the CPU stressor should execute the busy-wait loop on. Default: 1

#### Format

```json
"cpu_complexity": {
  "execution_time": <float:seconds>,
  "threads": <integer>
}
```

### Network Complexity

#### Optional attributes

* **forward_requests**: Determines if several calls to endpoints should be made in parallel. Default: "synchronous"
* **response_payload_size**: Determines the number of characters that the server should send back to the calling service or stressor. Default: 0
* **called_services**: An array of endpoints that this endpoint will call before responding. Default: empty array

```json
"network_complexity": {
  "forward_requests": "<string:synchronous|asynchronous>",
  "response_payload_size": <integer:chars>,

  "called_services": [...]
}
```

### Endpoint Calls

#### Required attributes

* **service**: The name of the service object that serves the specified endpoint.
* **endpoint**: The name of the endpoint that will be contacted.
* **port**: The port the server is responding to requests on.
* **protocol**: Determines if the call will be made using HTTP or gRPC.
* **traffic_forward_ratio**: Determines the ratio of inbound to outbound requests (1:X). This determines how many times an endpoint call will be made for every request.

#### Optional attributes

* **request_payload_size**: Determines the number of characters that will be sent in the request to the endpoint. Default: 0

#### Format

```json
"called_services": [
  {
    "service": "<string>",
    "endpoint": "<string>",
    "port": "<string>",
    "protocol": "<string:http|grpc>",
    "traffic_forward_ratio": <integer>,
    "request_payload_size": <integer:chars>
  }
]
```

# Examples

Examples for simple and complex applications generated with HydraGen can be found [here](https://github.com/EricssonResearch/cloud-native-app-simulator/tree/main/generator/examples). The .json is the taxonomy description give as input to the application generator and the the clusterX folder(s) contain the Kubernetes .yaml files generated by this module.
