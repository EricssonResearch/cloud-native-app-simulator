/*
Copyright 2023 Telefonaktiebolaget LM Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stressors

import (
	"application-emulator/src/util"
	model "application-model"
	"application-model/generated"
	"fmt"
	"math/rand"
	"strings"
)

// Characters in response payload
const characters = "abcdefghijklmnopqrstuvwxyz"

// Generates a random payload of size n
func RandomPayload(n int) string {
	if n == 0 {
		return ""
	}

	builder := strings.Builder{}
	builder.Grow(n)

	for i := 0; i < n; i++ {
		builder.WriteByte(characters[rand.Int()%len(characters)])
	}

	return builder.String()
}

type NetworkTask struct {
	Request any
}

// Combines the task responses in taskResponses with networkTaskResponse and endpointResponses
func ConcatenateNetworkResponses(taskResponses *MutexTaskResponses, networkTaskResponse *generated.NetworkTaskResponse, endpointResponses []generated.EndpointResponse) {
	taskResponses.Mutex.Lock()
	defer taskResponses.Mutex.Unlock()

	if taskResponses.NetworkTask != nil {
		taskResponses.NetworkTask.Services = append(taskResponses.NetworkTask.Services, networkTaskResponse.Services...)
		for k, v := range networkTaskResponse.Responses {
			taskResponses.NetworkTask.Responses[k] = v
		}
		// Don't replace the payload
	} else {
		taskResponses.NetworkTask = networkTaskResponse
	}

	for _, r := range endpointResponses {
		key := fmt.Sprintf("%s/%s", r.Service.Service, r.Service.Endpoint)
		taskResponses.NetworkTask.Responses[key] = &generated.ServiceResponse{
			Protocol: r.Protocol,
			Status:   r.Status,
		}

		// ResponseData is nil if an error occured
		// ResponseData.Tasks is nil if the endpoint was not found or ran no tasks
		if r.ResponseData != nil && r.ResponseData.Tasks != nil {
			taskResponses.Mutex.Unlock()
			if r.ResponseData.Tasks.CpuTask != nil {
				ConcatenateCPUResponses(taskResponses, r.ResponseData.Tasks.CpuTask)
			}
			if r.ResponseData.Tasks.NetworkTask != nil {
				ConcatenateNetworkResponses(taskResponses, r.ResponseData.Tasks.NetworkTask, nil)
			}
			taskResponses.Mutex.Lock()
		}
	}
}

func (n *NetworkTask) ExecAllowed(endpoint *model.Endpoint) bool {
	return endpoint.NetworkComplexity != nil
}

// Stress the network by returning a user-defined payload and calling other endpoints
func (n *NetworkTask) ExecTask(endpoint *model.Endpoint, responses *MutexTaskResponses) {
	stressParams := endpoint.NetworkComplexity

	var calls []generated.EndpointResponse
	if stressParams.ForwardRequests == "asynchronous" {
		calls = ForwardParallel(n.Request, stressParams.CalledServices)
	} else if stressParams.ForwardRequests == "synchronous" {
		calls = ForwardSequential(n.Request, stressParams.CalledServices)
	}

	svc := fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)
	ConcatenateNetworkResponses(responses, &generated.NetworkTaskResponse{
		Services:  []string{svc},
		Responses: make(map[string]*generated.ServiceResponse),
		Payload:   RandomPayload(stressParams.ResponsePayloadSize),
	}, calls)

	util.LogNetworkTask(endpoint, calls)
}
