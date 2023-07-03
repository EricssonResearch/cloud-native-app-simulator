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

type EndpointTrace struct {
	Endpoint *model.Endpoint
	Time     time.Time
	CPUTime  int64
}

var loggingEnabled = false

func SetLoggingEnabled(enabled bool) string {
	loggingEnabled = enabled

	if loggingEnabled {
		return "logging enabled"
	} else {
		return "logging disabled"
	}
}

// Call at start of endpoint call to trace execution time
func TraceEndpointStart(endpoint *model.Endpoint) *EndpointTrace {
	if loggingEnabled {
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
func TraceEndpointEnd(trace *EndpointTrace) {
	if trace != nil {
		responseTime := time.Now().Sub(trace.Time).Seconds()
		cpuTime := float64(ProcessCPUTime()-trace.CPUTime) / 1000000000.0

		log.Printf("%s %s: responseTime=%fs cpuTime=%fs",
			trace.Endpoint.Protocol, trace.Endpoint.Name, responseTime, cpuTime)
	}
}
