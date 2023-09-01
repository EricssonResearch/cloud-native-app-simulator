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
	"runtime"
	"sync"
)

type CPUTask struct{}

// Combines the CPU task response in taskResponses with cpuTaskResponse
func ConcatenateCPUResponses(taskResponses *MutexTaskResponses, cpuTaskResponse *generated.CPUTaskResponse) {
	taskResponses.Mutex.Lock()
	defer taskResponses.Mutex.Unlock()

	if taskResponses.CpuTask != nil {
		for k, v := range cpuTaskResponse.Services {
			uniqueKey := UniqueKey(taskResponses.CpuTask.Services, k)
			taskResponses.CpuTask.Services[uniqueKey] = v
		}
	} else {
		taskResponses.CpuTask = cpuTaskResponse
	}
}

func (c *CPUTask) ExecAllowed(endpoint *model.Endpoint) bool {
	return endpoint.CpuComplexity != nil
}

func StressCPU(executionTime float32, lockThread bool) {
	if executionTime > 0 {
		// Threads need to be locked because otherwise util.ThreadCPUTime() can change in the middle of execution
		if lockThread {
			runtime.LockOSThread()
		}

		start := util.ThreadCPUTime()
		target := start + int64(executionTime*1000000000.0)

		for util.ThreadCPUTime() < target {
		}

		if lockThread {
			runtime.UnlockOSThread()
		}
	}
}

// Stress the CPU by running a busy loop, if the endpoint has a defined CPU complexity
func (c *CPUTask) ExecTask(endpoint *model.Endpoint, responses *MutexTaskResponses) {
	stressParams := endpoint.CpuComplexity

	if stressParams.Threads > 1 {
		wg := sync.WaitGroup{}
		wg.Add(stressParams.Threads)

		for i := 0; i < stressParams.Threads; i++ {
			go func() {
				defer wg.Done()
				StressCPU(stressParams.ExecutionTime, true)
			}()
		}

		wg.Wait()
	} else {
		StressCPU(stressParams.ExecutionTime, true)
	}

	svc := fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)
	ConcatenateCPUResponses(responses, &generated.CPUTaskResponse{
		Services: map[string]float32{
			svc: stressParams.ExecutionTime,
		},
	})

	util.LogCPUTask(endpoint)
}
