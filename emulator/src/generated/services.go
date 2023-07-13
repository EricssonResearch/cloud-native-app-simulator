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
	"log"

	"google.golang.org/grpc"
)

// This file will be replaced by the generated gRPC code when the emulator executes

type Service1ServerImpl struct {
	UnimplementedService1Server
	Endpoint *model.Endpoint
}

func (s *Service1ServerImpl) TestEndpoint(ctx context.Context, request *generated.Request) (*generated.Response, error) {
	trace := util.TraceEndpointCall(s.Endpoint)
	response := &generated.Response{
		Endpoint: s.Endpoint.Name,
		Tasks:    stressors.Exec(request, s.Endpoint),
	}
	util.LogEndpointCall(trace)
	// TODO: When to return error?
	return response, nil
}

// Maps endpoint to a generated service struct and registers it with registrar
func RegisterGeneratedService(registrar grpc.ServiceRegistrar, endpoint *model.Endpoint) {
	switch endpoint.Name {
	case "service-1":
		RegisterService1Server(registrar, &Service1ServerImpl{Endpoint: endpoint})
	default:
		log.Fatalf("Tried to register gRPC service %s with no generated code", endpoint.Name)
	}
}
