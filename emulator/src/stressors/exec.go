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
	model "application-model"
	"sync"
)

// Interface for a stressor used to simulate the workload of a microservice
type Stressor interface {
	// If the stressor should execute according to the parameters provided by the user
	ExecAllowed(endpoint *model.Endpoint) bool
	// Executes the workload according to user parameters
	ExecTask(endpoint *model.Endpoint) any
	// Combine REST and gRPC responses from endpoint calls
	CombineResponses(selfResponse any, endpointResponse any)
}

// Recursively scan the result from the network stressor for endpoint responses to combine
func combineNetworkTaskResponses(allResponses map[string]any, stressors map[string]Stressor) {
	if networkTaskResponse, ok := allResponses["stress_network"]; ok {
		networkTaskResponse := networkTaskResponse.(*model.NetworkTaskResponse)

		for _, endpointResponse := range networkTaskResponse.EndpointResponses {
			if restResponse, ok := endpointResponse.(*model.RESTResponse); ok {
				combineNetworkTaskResponses(restResponse.Tasks, stressors)

				for name, stressor := range stressors {
					selfResponse, selfOk := allResponses[name]
					endpointResponse, endpointOk := restResponse.Tasks[name]

					// If both our response and the response from the endpoint contains this stressor, combine the result
					if selfOk && endpointOk {
						stressor.CombineResponses(selfResponse, endpointResponse)
					}
				}
			}
		}

		// Should not be included in response JSON
		networkTaskResponse.EndpointResponses = nil
	}
}

// Executes all stressors defined in the endpoint sequentially
func ExecSequential(request any, endpoint *model.Endpoint) map[string]any {
	cpuTask := &CPUTask{}
	networkTask := &NetworkTask{Request: request}
	stressors := map[string]Stressor{
		"stress_cpu":     cpuTask,
		"stress_network": networkTask,
	}
	responses := make(map[string]any)

	for name, stressor := range stressors {
		if stressor.ExecAllowed(endpoint) {
			responses[name] = stressor.ExecTask(endpoint)
		}
	}

	combineNetworkTaskResponses(responses, stressors)
	return responses
}

func execStressor(name string, stressor Stressor, endpoint *model.Endpoint, responses map[string]any, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	response := stressor.ExecTask(endpoint)
	mutex.Lock()
	responses[name] = response
	mutex.Unlock()
}

// Executes all stressors defined in the endpoint in parallel using goroutines
func ExecParallel(request any, endpoint *model.Endpoint) map[string]any {
	cpuTask := &CPUTask{}
	networkTask := &NetworkTask{Request: request}
	stressors := map[string]Stressor{
		"stress_cpu":     cpuTask,
		"stress_network": networkTask,
	}
	responses := make(map[string]any)

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for name, stressor := range stressors {
		if stressor.ExecAllowed(endpoint) {
			wg.Add(1)
			go execStressor(name, stressor, endpoint, responses, &wg, &mutex)
		}
	}

	wg.Wait()
	// Execute this sequentially since concurrent map access is not allowed
	combineNetworkTaskResponses(responses, stressors)

	return responses
}
