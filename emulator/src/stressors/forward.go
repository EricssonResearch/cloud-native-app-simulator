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
	"application-model/generated"
	"fmt"
	"log"
	"net/http"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// Override the content type
	forwardHeaders.Set("Content-Type", "application/json")

	return forwardHeaders
}

func httpRequest(service model.CalledService, forwardHeaders http.Header, sourceEndpoint string) generated.EndpointResponse {
	log.Printf("[FORWARD] Http request ready to send POST to %s/%s", service.Service, service.Endpoint)
	status, response, err :=
		client.POST(service.Service, service.Endpoint, service.Port, RandomPayload(service.RequestPayloadSize), forwardHeaders, sourceEndpoint)

	log.Printf("[FORWARD] Http request returned from %s/%s", service.Service, service.Endpoint)
	if err != nil {
		return generated.EndpointResponse{
			Service:  &service,
			Status:   err.Error(),
			Protocol: "HTTP",
		}
	} else {
		return generated.EndpointResponse{
			Service:      &service,
			Status:       fmt.Sprintf("%d %s", status, http.StatusText(status)),
			Protocol:     "HTTP",
			ResponseData: response,
		}
	}
}

func grpcRequest(service model.CalledService, sourceEndpoint string) generated.EndpointResponse {
	response, err :=
		client.GRPC(service.Service, service.Endpoint, service.Port, RandomPayload(service.RequestPayloadSize), sourceEndpoint)

	if err != nil {
		return generated.EndpointResponse{
			Service:  &service,
			Status:   status.Convert(err).Code().String(),
			Protocol: "gRPC",
		}
	} else {
		return generated.EndpointResponse{
			Service:      &service,
			Status:       codes.OK.String(),
			Protocol:     "gRPC",
			ResponseData: response,
		}
	}
}

// Forward requests to all services sequentially and return REST or gRPC responses
func ForwardSequential(request any, services []model.CalledService, endpoint string) []generated.EndpointResponse {
	forwardHeaders := ExtractHeaders(request)
	len := 0
	for _, service := range services {
		len += service.TrafficForwardRatio
	}
	responses := make([]generated.EndpointResponse, len, len)

	i := 0
	for _, service := range services {
		for j := 0; j < service.TrafficForwardRatio; j++ {
			if service.Protocol == "http" {
				response := httpRequest(service, forwardHeaders, endpoint)
				responses[i] = response
			} else if service.Protocol == "grpc" {
				response := grpcRequest(service, endpoint)
				responses[i] = response
			}
			i++
		}
	}

	return responses
}

func parallelHTTPRequest(responses []generated.EndpointResponse, i int, service model.CalledService, forwardHeaders http.Header, wg *sync.WaitGroup, endpoint string) {
	defer wg.Done()
	response := httpRequest(service, forwardHeaders, endpoint)
	// No mutex needed since every response has its own index
	responses[i] = response
}

func parallelGRPCRequest(responses []generated.EndpointResponse, i int, service model.CalledService, wg *sync.WaitGroup, endpoint string) {
	defer wg.Done()
	response := grpcRequest(service, endpoint)
	// No mutex needed since every response has its own index
	responses[i] = response
}

// Forward requests to all services in parallel using goroutines and return REST or gRPC responses
func ForwardParallel(request any, services []model.CalledService, endpoint string) []generated.EndpointResponse {
	forwardHeaders := ExtractHeaders(request)
	len := 0
	for _, service := range services {
		len += service.TrafficForwardRatio
	}
	responses := make([]generated.EndpointResponse, len, len)
	wg := sync.WaitGroup{}

	i := 0
	for _, service := range services {
		for j := 0; j < service.TrafficForwardRatio; j++ {
			if service.Protocol == "http" {
				wg.Add(1)
				go parallelHTTPRequest(responses, i, service, forwardHeaders, &wg, endpoint)
			} else if service.Protocol == "grpc" {
				wg.Add(1)
				go parallelGRPCRequest(responses, i, service, &wg, endpoint)
			}
			i++
		}
	}

	wg.Wait()
	return responses
}
