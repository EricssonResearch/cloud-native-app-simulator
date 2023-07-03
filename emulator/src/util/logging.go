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

package util

import (
	"cloud-native-app-simulator/model"

	"log"
	"time"
)

// TODO: Move these
var ServiceName = ""
var LoggingEnabled = false

type EndpointTrace struct {
	Endpoint *model.Endpoint
	Time     time.Time
	CPUTime  int64
}

// Call at start of endpoint call to trace execution time
func TraceEndpointCall(endpoint *model.Endpoint) *EndpointTrace {
	if LoggingEnabled {
		trace := &EndpointTrace{
			Endpoint: endpoint,
			Time:     time.Now(),
			CPUTime:  ProcessCPUTime(),
		}

		return trace
	} else {
		return nil
	}
}

// Call at end of endpoint call to print stats to stdout
func PrintEndpointTrace(trace *EndpointTrace) {
	if trace != nil {
		responseTime := time.Now().Sub(trace.Time).Seconds()
		cpuTime := float64(ProcessCPUTime()-trace.CPUTime) / 1000000000.0

		log.Printf("%s %s/%s: %s responseTime=%fs cpuTime=%fs",
			trace.Endpoint.Protocol, ServiceName, trace.Endpoint.Name, trace.Endpoint.ExecutionMode, responseTime, cpuTime)
	}
}
