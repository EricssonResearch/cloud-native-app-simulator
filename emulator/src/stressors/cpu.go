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

	"runtime"
)

type CPUTaskResponse struct {
	Services []string `json:"services"`
	Statuses []string `json:"statuses"`
}

// Stress the CPU by running a busy loop, if the endpoint has a defined CPU complexity
func CPU(cpuComplexity *model.CpuComplexity) {
	// TODO: This needs to be tested more
	if executionTime := cpuComplexity.ExecutionTime; executionTime > 0 {
		runtime.LockOSThread()

		start := util.ThreadCPUTime()
		target := start + int64(executionTime)*1000000000

		for util.ThreadCPUTime() < target {
		}

		runtime.UnlockOSThread()
	}
}
