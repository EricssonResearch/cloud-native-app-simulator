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

// This file will be replaced by the generated gRPC code when the emulator image is built

import (
	generated "application-emulator/src/generated"
	generated_model "application-model/generated"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Searches for method by service, endpoint and returns the result
func CallGeneratedEndpoint(ctx context.Context, cc grpc.ClientConnInterface, service, endpoint string, in *generated_model.Request, options ...grpc.CallOption) (*generated_model.Response, error) {
	switch service {
	case "service-1":
		client := generated.NewService1Client(cc)
		switch endpoint {
		case "test-endpoint":
			return client.TestEndpoint(ctx, in, options...)
		}
	}

	return nil, status.Error(codes.InvalidArgument, "unknown service, endpoint combination")
}
