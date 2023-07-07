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

package model

type CPUTaskResponse struct {
	Services []string `json:"services"`
	Statuses []string `json:"statuses"`
}

type NetworkTaskResponse struct {
	Services []string `json:"services"`
	Statuses []string `json:"statuses"`
	Payload  string   `json:"payload"`
}

type TaskResponses struct {
	CPUTask     *CPUTaskResponse     `json:"cpu_task,omitempty"`
	NetworkTask *NetworkTaskResponse `json:"network_task,omitempty"`
}

type RESTRequest struct {
	Payload string `json:"payload,omitempty"`
}

type RESTResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"message,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`

	Tasks TaskResponses `json:"tasks"`
}

type GRPCResponse struct {
	// TODO
}
