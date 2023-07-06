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
	"runtime"
)

type CPUTask struct{}

func (c *CPUTask) ExecAllowed(endpoint *model.Endpoint) bool {
	return endpoint.CpuComplexity != nil
}

// Stress the CPU by running a busy loop, if the endpoint has a defined CPU complexity
func (c *CPUTask) ExecTask(endpoint *model.Endpoint) any {
	stressParams := endpoint.CpuComplexity

	// TODO: This needs to be tested more
	if stressParams.ExecutionTime > 0 {
		runtime.LockOSThread()

		start := util.ThreadCPUTime()
		target := start + int64(stressParams.ExecutionTime)*1000000000

		for util.ThreadCPUTime() < target {
		}

		runtime.UnlockOSThread()
	}

	return &model.CPUTaskResponse{
		Services: []string{fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)},
		Statuses: []string{fmt.Sprintf("execution_time: %f", stressParams.ExecutionTime)},
	}
}

func (c *CPUTask) CombineResponses(s any, e any) {
	selfResponse := s.(*model.CPUTaskResponse)
	endpointResponse := e.(*model.CPUTaskResponse)

	selfResponse.Services = append(endpointResponse.Services, selfResponse.Services...)
	selfResponse.Statuses = append(endpointResponse.Statuses, selfResponse.Statuses...)
}
