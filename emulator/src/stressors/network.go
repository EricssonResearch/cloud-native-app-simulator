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

	"fmt"
	"math/rand"
	"strings"
)

type NetworkTask struct {
	Request any
}

// Characters in response payload
const characters = "abcdefghijklmnopqrstuvwxyz"

// Generates a random payload of size n
func RandomPayload(n int) string {
	builder := strings.Builder{}
	builder.Grow(n)

	for i := 0; i < n; i++ {
		builder.WriteByte(characters[rand.Int()%len(characters)])
	}

	return builder.String()
}

// Combines the task responses in taskResponses with networkTaskResponse and endpointResponses
func ConcatenateNetworkResponses(taskResponses *MutexTaskResponses, networkTaskResponse *model.NetworkTaskResponse, endpointResponses []EndpointResponse) {
	taskResponses.Mutex.Lock()
	defer taskResponses.Mutex.Unlock()

	if networkTaskResponse == nil {
		return
	}

	if taskResponses.NetworkTask != nil {
		taskResponses.NetworkTask.Services = append(networkTaskResponse.Services, taskResponses.NetworkTask.Services...)
		taskResponses.NetworkTask.Statuses = append(networkTaskResponse.Statuses, taskResponses.NetworkTask.Statuses...)
		// Don't replace the payload
	} else {
		taskResponses.NetworkTask = networkTaskResponse
	}

	for _, r := range endpointResponses {
		taskResponses.NetworkTask.Services = append([]string{r.Service}, taskResponses.NetworkTask.Services...)
		taskResponses.NetworkTask.Statuses = append([]string{r.Status}, taskResponses.NetworkTask.Statuses...)

		if rr := r.RESTResponse; rr != nil {
			taskResponses.Mutex.Unlock()
			ConcatenateCPUResponses(taskResponses, rr.Tasks.CPUTask)
			ConcatenateNetworkResponses(taskResponses, rr.Tasks.NetworkTask, nil)
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

	var calls []EndpointResponse
	if stressParams.ForwardRequests == "asynchronous" {
		calls = ForwardParallel(n.Request, stressParams.CalledServices)
	} else if stressParams.ForwardRequests == "synchronous" {
		calls = ForwardSequential(n.Request, stressParams.CalledServices)
	}

	services := make([]string, len(calls)+1)
	services = append(services, fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name))
	statuses := make([]string, len(calls))

	for _, endpoint := range calls {
		services = append(services, endpoint.Service)
		statuses = append(statuses, endpoint.Status)
	}

	ConcatenateNetworkResponses(responses, &model.NetworkTaskResponse{
		Services: services,
		Statuses: statuses,
		Payload:  RandomPayload(stressParams.ResponsePayloadSize),
	}, calls)
}
