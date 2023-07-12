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
	generated "application-model/generated"
	"sync"
)

// Stressors need to access model.TaskResponses in parallel
type MutexTaskResponses struct {
	*sync.Mutex
	*generated.TaskResponses
}

// Interface for a stressor used to simulate the workload of a microservice
type Stressor interface {
	// If the stressor should execute according to the parameters provided by the user
	ExecAllowed(endpoint *model.Endpoint) bool
	// Executes the workload according to user parameters
	ExecTask(endpoint *model.Endpoint, responses *MutexTaskResponses)
}

// Executes all stressors defined in the endpoint sequentially
func ExecSequential(request any, endpoint *model.Endpoint) *generated.TaskResponses {
	stressors := []Stressor{
		&CPUTask{},
		&NetworkTask{Request: request},
	}
	responses := MutexTaskResponses{
		&sync.Mutex{},
		&generated.TaskResponses{},
	}

	for _, stressor := range stressors {
		if stressor.ExecAllowed(endpoint) {
			stressor.ExecTask(endpoint, &responses)
		}
	}

	return responses.TaskResponses
}

func execStressor(stressor Stressor, endpoint *model.Endpoint, responses *MutexTaskResponses, wg *sync.WaitGroup) {
	defer wg.Done()
	// responses can be accessed in parallel as long as
	stressor.ExecTask(endpoint, responses)
}

// Executes all stressors defined in the endpoint in parallel using goroutines
func ExecParallel(request any, endpoint *model.Endpoint) *generated.TaskResponses {
	stressors := []Stressor{
		&CPUTask{},
		&NetworkTask{Request: request},
	}
	responses := MutexTaskResponses{
		&sync.Mutex{},
		&generated.TaskResponses{},
	}
	wg := sync.WaitGroup{}

	for _, stressor := range stressors {
		if stressor.ExecAllowed(endpoint) {
			wg.Add(1)
			go execStressor(stressor, endpoint, &responses, &wg)
		}
	}

	wg.Wait()
	return responses.TaskResponses
}
