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
	"sync"
)

// Executes all stressors defined in the endpoint sequentially
func Exec(endpoint *model.Endpoint) *CPUTaskResponse {
	var cpuTaskResponse *CPUTaskResponse

	if endpoint.CpuComplexity != nil {
		CPU(endpoint.CpuComplexity)

		service := fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)
		// TODO: More information could be provided here
		executionTime := fmt.Sprintf("execution_time: %f", endpoint.CpuComplexity.ExecutionTime)

		cpuTaskResponse = &CPUTaskResponse{
			Services: []string{service},
			Statuses: []string{executionTime},
		}
	}

	return cpuTaskResponse
}

// Executes all stressors defined in the endpoint in parallel using goroutines
func ExecParallel(endpoint *model.Endpoint) *CPUTaskResponse {
	var cpuTaskResponse *CPUTaskResponse
	wg := sync.WaitGroup{}

	if endpoint.CpuComplexity != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()
			CPU(endpoint.CpuComplexity)

			service := fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)
			executionTime := fmt.Sprintf("execution_time: %f", endpoint.CpuComplexity.ExecutionTime)

			cpuTaskResponse = &CPUTaskResponse{
				Services: []string{service},
				Statuses: []string{executionTime},
			}
		}()
	}

	wg.Wait()
	return cpuTaskResponse
}
