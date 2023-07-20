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

import (
	"application-emulator/src/util"
	model "application-model"
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type HealthServerImpl struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (h *HealthServerImpl) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	grpcServiceName := fmt.Sprintf("generated.%s", model.GoName(util.ServiceName))
	if request.Service == "" || request.Service == grpcServiceName {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	} else {
		// https://github.com/grpc/grpc/blob/master/src/proto/grpc/health/v1/health.proto#L44
		return nil, status.Errorf(codes.NotFound, "Only serving %s", grpcServiceName)
	}
}

// Launch a gRPC server to serve one or more endpoints
func GRPC(endpoints []model.Endpoint) {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	util.GRPCRegisterGeneratedService(server, endpoints)
	reflection.Register(server)
	grpc_health_v1.RegisterHealthServer(server, &HealthServerImpl{})

	err = server.Serve(listener)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		panic(err)
	}
}
