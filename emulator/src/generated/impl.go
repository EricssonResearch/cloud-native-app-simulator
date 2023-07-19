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

package generated

import (
	"application-emulator/src/stressors"
	"application-emulator/src/util"
	model "application-model"
	"application-model/generated"
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
)

// This file will be replaced by the generated gRPC code when the emulator executes

type Service1ServerImpl struct {
	UnimplementedService1Server
	TestEndpointInfo *model.Endpoint
}

func (s *Service1ServerImpl) TestEndpoint(ctx context.Context, request *generated.Request) (*generated.Response, error) {
	trace := util.TraceEndpointCall(s.TestEndpointInfo, "gRPC")
	response := &generated.Response{
		Endpoint: s.TestEndpointInfo.Name,
		Tasks:    stressors.Exec(request, s.TestEndpointInfo),
	}
	util.LogEndpointCall(trace)
	return response, nil
}

// Maps endpoint to a generated service struct and registers it with registrar
func RegisterGeneratedService(registrar grpc.ServiceRegistrar, endpoints []model.Endpoint) {
	// Service name is empty in test setup
	switch util.ServiceName {
	case "service-1":
		impl := Service1ServerImpl{}
		for _, endpoint := range endpoints {
			switch endpoint.Name {
			case "test-endpoint":
				impl.TestEndpointInfo = &endpoint
			default:
				log.Fatalf("Service %s got invalid gRPC endpoint %s", util.ServiceName, endpoint.Name)
			}
		}
		RegisterService1Server(registrar, &impl)
	}
}

// Searches for method by service, endpoint and returns the result
func CallGeneratedEndpoint(ctx context.Context, cc grpc.ClientConnInterface, service, endpoint string, in *generated.Request) (*generated.Response, error) {
	options := []grpc.CallOption{}

	switch service {
	case "service-1":
		client := NewService1Client(cc)
		switch endpoint {
		case "test-endpoint":
			return client.TestEndpoint(ctx, in, options...)
		}
	}

	return nil, errors.New("Unknown service, endpoint combination")
}
