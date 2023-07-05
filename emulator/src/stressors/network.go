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
	"application-emulator/src/util"
	model "application-model"

	"fmt"
	"math/rand"
	"strings"
)

type NetworkTask struct {
	Request any
}

// Characters in response payload
const characters = "abcdefghijklmnopqrstuvwxyz"

// Generates a random payload of size n
func RandomPayload(n int) string {
	builder := strings.Builder{}
	builder.Grow(n)

	for i := 0; i < n; i++ {
		builder.WriteByte(characters[rand.Int()%len(characters)])
	}

	return builder.String()
}

func (n *NetworkTask) ExecAllowed(endpoint *model.Endpoint) bool {
	return endpoint.NetworkComplexity != nil
}

// Stress the network by returning a user-defined payload and calling other endpoints
func (n *NetworkTask) ExecTask(endpoint *model.Endpoint) any {
	stressParams := endpoint.NetworkComplexity
	responsePayload := RandomPayload(stressParams.ResponsePayloadSize)
	calls := []EndpointCall{}

	if stressParams.ForwardRequests == "asynchronous" {
		calls = ForwardParallel(n.Request, stressParams.CalledServices)
	} else if stressParams.ForwardRequests == "synchronous" {
		calls = ForwardSequential(n.Request, stressParams.CalledServices)
	}

	statuses := make([]string, 0, len(calls))
	for _, call := range calls {
		statuses = append(statuses, call.Status)
	}

	// TODO: Merge task responses
	return model.NetworkTaskResponse{
		Services: []string{fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)},
		Statuses: statuses,
		Payload:  responsePayload,
	}
}
