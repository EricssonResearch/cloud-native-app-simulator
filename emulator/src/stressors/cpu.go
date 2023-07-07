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

// Combines the CPU task response in taskResponses with cpuTaskResponse
func ConcatenateCPUResponses(taskResponses *MutexTaskResponses, cpuTaskResponse *model.CPUTaskResponse) {
	taskResponses.Mutex.Lock()
	defer taskResponses.Mutex.Unlock()

	if cpuTaskResponse == nil {
		return
	}

	if taskResponses.CPUTask != nil {
		taskResponses.CPUTask.Services = append(cpuTaskResponse.Services, taskResponses.CPUTask.Services...)
		taskResponses.CPUTask.Statuses = append(cpuTaskResponse.Statuses, taskResponses.CPUTask.Statuses...)
	} else {
		taskResponses.CPUTask = cpuTaskResponse
	}
}

func (c *CPUTask) ExecAllowed(endpoint *model.Endpoint) bool {
	return endpoint.CpuComplexity != nil
}

// Stress the CPU by running a busy loop, if the endpoint has a defined CPU complexity
func (c *CPUTask) ExecTask(endpoint *model.Endpoint, responses *MutexTaskResponses) {
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

	ConcatenateCPUResponses(responses, &model.CPUTaskResponse{
		Services: []string{fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)},
		Statuses: []string{fmt.Sprintf("execution_time: %f", stressParams.ExecutionTime)},
	})
}
