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
	"application-emulator/src/client"
	model "application-model"
	"fmt"
	"net/http"
	"sync"
)

// Headers to propagate from inbound to outbound
var incomingHeaders = []string{
	"User-Agent", "End-User", "X-Request-Id", "X-B3-TraceId", "X-B3-SpanId", "X-B3-ParentSpanId", "X-B3-Sampled", "X-B3-Flags",
}

// Extract relevant headers from the source request
func ExtractHeaders(request any) http.Header {
	// If this is a HTTP request, we should propagate the headers specified in incomingHeaders
	httpRequest, ok := request.(*http.Request)
	forwardHeaders := make(http.Header)

	if ok {
		for _, key := range incomingHeaders {
			if value := httpRequest.Header.Get(key); value != "" {
				forwardHeaders.Set(key, value)
			}
		}
	}

	return forwardHeaders
}

type EndpointResponse struct {
	Status       string
	RESTResponse *model.RESTResponse
	GRPCResponse *model.GRPCResponse
}

func httpRequest(service model.CalledService, forwardHeaders http.Header) EndpointResponse {
	status, response, err :=
		client.POST(service.Service, service.Endpoint, service.Port, RandomPayload(service.RequestPayloadSize), forwardHeaders)

	if err != nil {
		return EndpointResponse{
			Status: err.Error(),
		}
	} else {
		return EndpointResponse{
			// TODO: Should also return RESTResponse.Status?
			Status:       fmt.Sprintf("%d %s", status, http.StatusText(status)),
			RESTResponse: response,
		}
	}
}

// Forward requests to all services sequentially and return REST or gRPC responses
func ForwardSequential(request any, services []model.CalledService) []EndpointResponse {
	forwardHeaders := ExtractHeaders(request)
	responses := make([]EndpointResponse, len(services), len(services))

	for i, service := range services {
		// TODO: gRPC
		if service.Protocol == "http" {
			response := httpRequest(service, forwardHeaders)
			responses[i] = response
		}
	}

	return responses
}

func parallelHTTPRequest(responses []EndpointResponse, i int, service model.CalledService, forwardHeaders http.Header, wg *sync.WaitGroup) {
	defer wg.Done()
	response := httpRequest(service, forwardHeaders)
	// No mutex needed since every response has its own index
	responses[i] = response
}

// Forward requests to all services in parallel using goroutines and return REST or gRPC responses
func ForwardParallel(request any, services []model.CalledService) []EndpointResponse {
	forwardHeaders := ExtractHeaders(request)
	responses := make([]EndpointResponse, len(services), len(services))
	wg := sync.WaitGroup{}

	for i, service := range services {
		// TODO: gRPC
		if service.Protocol == "http" {
			wg.Add(1)
			go parallelHTTPRequest(responses, i, service, forwardHeaders, &wg)
		}
	}

	wg.Wait()
	return responses
}
