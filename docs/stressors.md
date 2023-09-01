# Adding a new stressor to the application emulator

If necessary, new stressors can be added to the application emulator to replicate the computational complexity or stress a microservice exerts on different hardware resources.

## Adding configuration parameters to the input file

The first step for adding a new stressor is to add it to the configmap that the emulator reads when starting.
The structure that contains endpoint configuration data is located in model/config.go.

First, add a new structure for the stressor:

```go
type CpuComplexity struct {
    ExecutionTime float32 `json:"execution_time"`
    Threads       int     `json:"threads"`
}

type NetworkComplexity struct {
    ForwardRequests     string          `json:"forward_requests"`
    ResponsePayloadSize int             `json:"response_payload_size"`
    CalledServices      []CalledService `json:"called_services"`
}

type MyStressorComplexity struct {
    MyVariable int `json:"my_variable"`
}
```

Then, add the structure to the endpoint configuration:

```go
type Endpoint struct {
    Name                 string                `json:"name"`
    ExecutionMode        string                `json:"execution_mode"`
    CpuComplexity        *CpuComplexity        `json:"cpu_complexity,omitempty"`
    NetworkComplexity    *NetworkComplexity    `json:"network_complexity,omitempty"`
    MyStressorComplexity *MyStressorComplexity `json:"my_stressor_complexity,omitempty"`
}
```

The stressor can now be added to an endpoint in the input JSON file:

```json
"endpoints": [
    {
        "name": "end1",
        "execution_mode": "sequential",
        "cpu_complexity": {
            "execution_time": 0.001
        },
        "my_stressor_complexity": {
            "my_variable": 2
        }
    },
]
```

## Implementing the stressor

Stressors should be added in the directory emulator/src/stressors and must implement two functions:

```go
package stressors

type MyStressorTask struct{}

// Determines if the stressor should execute according to the parameters provided by the user
func (m *MyStressorTask) ExecAllowed(endpoint *model.Endpoint) bool { ... }
// Executes the workload according to user parameters
func (m *MyStressorTask) ExecTask(endpoint *model.Endpoint, responses *MutexTaskResponses) { ... }
```

The stressor should add a response to the task responses structure, which is located in model/api.proto:

```go
message CPUTaskResponse {
    map<string, float> services = 1;
}

message ServiceResponse {
    string protocol = 1;
    string status = 2;
}

message NetworkTaskResponse {
    repeated string services = 1;
    map<string, ServiceResponse> responses = 2;
    
    string payload = 3;
}

message MyTaskResponse {
    map<string, int> services = 1;
}

message TaskResponses {
    CPUTaskResponse cpu_task = 1;
    NetworkTaskResponse network_task = 2;
    MyTaskResponse my_task = 3;
}
```

An example implementation of a stressor that sleeps for the number of seconds specified by the user looks like this:

```go
func (m *MyStressorTask) ExecAllowed(endpoint *model.Endpoint) bool {
    return endpoint.MyStressorComplexity != nil
}

func (m *MyStressorTask) ExecTask(endpoint *model.Endpoint, responses *MutexTaskResponses) {
    stressParams := endpoint.MyStressorComplexity
    time.Sleep(stressParams.MyVariable * time.Second)

    svc := fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)
    ConcatenateMyStressorResponses(responses, &generated.MyTaskResponse{
        Services: map[string]int{
            svc: stressParams.MyVariable,
        },
    })

    util.LogMyTask(endpoint)
}
```

Logging for stressors should be added in util/logging.go:

```go
// Call at end of "my task" to print params to stdout
func LogMyTask(endpoint *model.Endpoint) {
    if LoggingEnabled {
        myVariable := endpoint.MyStressorComplexity.MyVariable
        log.Printf("%s/%s: My task myVariable=%d",
            ServiceName, endpoint.Name, myVariable)
    }
}
```

## Concatenating responses

The network stressor will concatenate responses it receives from other endpoints, which means our new stressor needs a `ConcatenateMyStressorResponses`. The function should append the response to the current list of responses.

```go
func ConcatenateMyStressorResponses(taskResponses *MutexTaskResponses, myTaskResponse *generated.MyTaskResponse) {
    taskResponses.Mutex.Lock()
    defer taskResponses.Mutex.Unlock()

    if taskResponses.MyTask != nil {
        for k, v := range myTaskResponse.Services {
            taskResponses.MyTask.Services[k] = v
        }
    } else {
        taskResponses.MyTask = myTaskResponse
    }
}
```

Add the new function to the network stressor in emulator/src/stressors/network.go:

```go
    for _, r := range endpointResponses {
        key := fmt.Sprintf("%s/%s", r.Service.Service, r.Service.Endpoint)
        taskResponses.NetworkTask.Responses[key] = &generated.ServiceResponse{
            Protocol: r.Protocol,
            Status:   r.Status,
        }

        if r.ResponseData != nil && r.ResponseData.Tasks != nil {
            taskResponses.Mutex.Unlock()
            if r.ResponseData.Tasks.CpuTask != nil {
                ConcatenateCPUResponses(taskResponses, r.ResponseData.Tasks.CpuTask)
            }
            if r.ResponseData.Tasks.NetworkTask != nil {
                ConcatenateNetworkResponses(taskResponses, r.ResponseData.Tasks.NetworkTask, nil)
            }
            if r.ResponseData.Tasks.MyTask != nil {
                ConcatenateMyStressorResponses(taskResponses, r.ResponseData.Tasks.MyTask)
            }
            taskResponses.Mutex.Lock()
        }
    }
```

## Executing the stressor when a request is received

The stressor needs to added to the list in both `ExecSequential` and `ExecParallel` to be executed when a request is received:

```go
stressors := []Stressor{
    &CPUTask{},
    &NetworkTask{Request: request},
    &MyTask{}
}
```

Done!
