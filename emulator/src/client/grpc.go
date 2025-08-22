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

package client

import (
	"application-emulator/src/generated/client"
	"application-emulator/src/resilience/circuit_breaker"
	"application-model/generated"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Sends a gRPC request to the specified endpoint
func GRPC(service, endpoint string, port int, payload, sourceEndpoint string) (*generated.Response, error) {
	var url string
	// Omit the port if zero
	if port == 0 {
		url = service
	} else {
		url = fmt.Sprintf("%s:%d", service, port)
	}

	// TODO: TLS?
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(url, dialOptions...)

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	circuitBreakerRegistry := circuit_breaker.GetCircuitBreakerRegistry()
	cbName := circuitBreakerRegistry.BuildName(sourceEndpoint, service, endpoint)
	circuitBreaker := circuitBreakerRegistry.GetCircuitBreaker(cbName)

	var response *generated.Response
	request := &generated.Request{
		Payload: payload,
	}
	callOptions := []grpc.CallOption{}

	if circuitBreaker == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err = client.CallGeneratedEndpoint(ctx, conn, service, endpoint, request, callOptions...)
	} else {
		response, err = circuitBreaker.ProxyGRPC(conn, service, endpoint, request, callOptions...)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}
