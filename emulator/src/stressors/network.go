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

type NetworkTask struct{}
type NetworkTaskResponse struct {
	Services []string `json:"services"`
	Statuses []string `json:"statuses"`
	Payload  string   `json:"payload"`
}

const characters = "abcdefghijklmnopqrstuvwxyz"

func (n *NetworkTask) ShouldExec(endpoint *model.Endpoint) bool {
	return endpoint.NetworkComplexity != nil
}

// Stress the network by returning a user-defined payload and calling other endpoints
func (n *NetworkTask) ExecTask(endpoint *model.Endpoint) any {
	payloadSize := endpoint.NetworkComplexity.ResponsePayloadSize

	builder := strings.Builder{}
	builder.Grow(payloadSize)

	for i := 0; i < payloadSize; i++ {
		builder.WriteByte(characters[rand.Int()%len(characters)])
	}

	return NetworkTaskResponse{
		Services: []string{fmt.Sprintf("%s/%s", util.ServiceName, endpoint.Name)},
		// TODO: Call other endpoints
		Statuses: []string{""},
		Payload:  builder.String(),
	}
}
