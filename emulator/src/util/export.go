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
	"context"

	"google.golang.org/grpc"
)

// This is necessary to break an import cycle between generated -> stressors -> client -> ...
type GRPCRegisterGeneratedServiceType = func(grpc.ServiceRegistrar, []model.Endpoint)
type GRPCCallGeneratedEndpointType = func(context.Context, grpc.ClientConnInterface, string, string, *generated.Request, ...grpc.CallOption) (*generated.Response, error)

var GRPCRegisterGeneratedService GRPCRegisterGeneratedServiceType
var GRPCCallGeneratedEndpoint GRPCCallGeneratedEndpointType
