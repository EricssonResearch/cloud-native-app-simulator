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
	"net/http"
	"strings"
)

type NetworkTask struct {
	Request any
}
type NetworkTaskResponse struct {
	Services []string `json:"services"`
	Statuses []string `json:"statuses"`
	Payload  string   `json:"payload"`
}

// Characters in response payload
const characters = "abcdefghijklmnopqrstuvwxyz"

// Headers to propagate from inbound to outbound
var incomingHeaders = []string{
	"user-agent", "end-user", "x-request-id", "x-b3-traceid", "x-b3-spanid", "x-b3-parentspanid", "x-b3-sampled", "x-b3-flags",
}

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

	// If this is a HTTP request, we should propagate the headers specified in incomingHeaders
	httpRequest, ok := n.Request.(*http.Request)
	forwardHeaders := make(http.Header)

	if ok {
		for _, key := range incomingHeaders {
			if value := httpRequest.Header.Get(key); value != "" {
				forwardHeaders.Set(key, value)
			}
		}
	}

	responsePayload := RandomPayload(stressParams.ResponsePayloadSize)
	return NetworkTaskResponse{
		Services: []string{fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)},
		// TODO: Call other endpoints
		Statuses: []string{""},
		Payload:  responsePayload,
	}
}
