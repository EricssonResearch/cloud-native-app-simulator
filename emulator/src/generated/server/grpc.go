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

package server

// This file will be replaced by the generated gRPC code when the emulator image is built

import (
	generated "application-emulator/src/generated"
	"application-emulator/src/stressors"
	"application-emulator/src/util"
	model "application-model"
	generated_model "application-model/generated"
	"context"
	"log"

	"google.golang.org/grpc"
)

type Service1ServerImpl struct {
	generated.UnimplementedService1Server
	TestEndpointInfo *model.Endpoint
}

func (s *Service1ServerImpl) TestEndpoint(ctx context.Context, request *generated_model.Request) (*generated_model.Response, error) {
	trace := util.TraceEndpointCall(s.TestEndpointInfo, "gRPC")
	response := &generated_model.Response{
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
		generated.RegisterService1Server(registrar, &impl)
	}
}
