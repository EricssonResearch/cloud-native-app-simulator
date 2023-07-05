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
	ShouldExec(endpoint *model.Endpoint) bool
	// Executes the workload according to user parameters
	ExecTask(endpoint *model.Endpoint) any
	// TODO: Combine responses from network complexity
}

var stressors = map[string]Stressor{
	"stress_cpu":     &CPUTask{},
	"stress_network": &NetworkTask{},
}

// Executes all stressors defined in the endpoint sequentially
func ExecSequential(endpoint *model.Endpoint) map[string]any {
	responses := make(map[string]any)

	for name, stressor := range stressors {
		if stressor.ShouldExec(endpoint) {
			responses[name] = stressor.ExecTask(endpoint)
		}
	}

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
func ExecParallel(endpoint *model.Endpoint) map[string]any {
	responses := make(map[string]any)
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for name, stressor := range stressors {
		if stressor.ShouldExec(endpoint) {
			wg.Add(1)
			go execStressor(name, stressor, endpoint, responses, &wg, &mutex)
		}
	}

	wg.Wait()
	return responses
}
