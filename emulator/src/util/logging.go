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
	model "application-model"
	"application-model/generated"
	"fmt"
	"log"
	"runtime"
	"time"
)

var ServiceName = "service-1"
var LoggingEnabled = false

type EndpointTrace struct {
	Endpoint *model.Endpoint
	Time     time.Time
	CPUTime  int64
}

// Call at start of program to print configuration to stdout
func LogConfiguration(configMap *model.ConfigMap) {
	// Get the process count from Go to make sure settings were applied
	processes := runtime.GOMAXPROCS(0)
	log.Printf("Application emulator started at *:5000-5001, logging: %t, processes: %d", LoggingEnabled, processes)

	httpEndpoints := []string{}
	grpcEndpoints := []string{}

	for _, endpoint := range configMap.Endpoints {
		if endpoint.Protocol == "http" {
			httpEndpoints = append(httpEndpoints, endpoint.Name)
		} else if endpoint.Protocol == "grpc" {
			grpcEndpoints = append(grpcEndpoints, fmt.Sprintf("generated.%s.%s", model.GoName(ServiceName), model.GoName(endpoint.Name)))
		}
	}

	log.Printf("Service: %s", ServiceName)
	log.Printf("HTTP endpoints: %v", httpEndpoints)
	log.Printf("gRPC endpoints: %v", grpcEndpoints)
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
func LogEndpointCall(trace *EndpointTrace) {
	if trace != nil {
		responseTime := time.Now().Sub(trace.Time).Seconds()
		cpuTime := float64(ProcessCPUTime()-trace.CPUTime) / 1000000000.0
		responseTimeFmt, cpuTimeFmt := FormatTime(responseTime), FormatTime(cpuTime)

		log.Printf("%s/%s: %s %s responseTime=%s cpuTime=%s",
			ServiceName, trace.Endpoint.Name, trace.Endpoint.Protocol, trace.Endpoint.ExecutionMode, responseTimeFmt, cpuTimeFmt)
	}
}

// Call at end of CPU task to print params to stdout
func LogCPUTask(endpoint *model.Endpoint) {
	if LoggingEnabled {
		executionTime := FormatTime(float64(endpoint.CpuComplexity.ExecutionTime))
		threads := endpoint.CpuComplexity.Threads

		log.Printf("%s/%s: CPU task executionTime=%s threads=%d lockThreads=%t",
			ServiceName, endpoint.Name, executionTime, threads, true)
	}
}

// Call at end of network task to print params to stdout
func LogNetworkTask(endpoint *model.Endpoint, responses []generated.EndpointResponse) {
	if LoggingEnabled {
		executionMode := endpoint.NetworkComplexity.ForwardRequests
		payloadSize := endpoint.NetworkComplexity.ResponsePayloadSize
		calledServices := len(endpoint.NetworkComplexity.CalledServices)

		statuses := make([]string, 0, len(responses))
		for _, response := range responses {
			statuses = append(statuses, fmt.Sprintf("http/%s:%s", response.ResponseData.Endpoint, response.Status))
		}
		formattedStatuses := fmt.Sprint(statuses)

		log.Printf("%s/%s: Network task %s payloadSize=%d calledServices=%d statuses=%s",
			ServiceName, endpoint.Name, executionMode, payloadSize, calledServices, formattedStatuses)
	}
}
