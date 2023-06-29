/*
Copyright 2021 Telefonaktiebolaget LM Ericsson AB

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

type GatewayInstance struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       GatewaySpec `yaml:"spec"`
}
type GatewaySelector struct {
	Istio string `yaml:"istio"`
}
type Port struct {
	Number   int    `yaml:"number,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}
type GatewayServers struct {
	Port  Port     `yaml:"port"`
	Hosts []string `yaml:"hosts"`
}
type GatewaySpec struct {
	Selector GatewaySelector  `yaml:"selector"`
	Servers  []GatewayServers `yaml:"servers"`
}

type VirtualServiceInstance struct {
	APIVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   Metadata           `yaml:"metadata"`
	Spec       VirtualServiceSpec `yaml:"spec"`
}
type Metadata struct {
	Name string `yaml:"name"`
}
type VirtualServiceURI struct {
	Exact string `yaml:"exact"`
}
type VirtualServiceMatch struct {
	URI VirtualServiceURI `yaml:"uri"`
}
type VirtualServiceDestination struct {
	Host string `yaml:"host"`
	Port Port   `yaml:"port"`
}
type VirtualServiceRoute struct {
	Destination VirtualServiceDestination `yaml:"destination"`
}
type VirtualServiceHTTP struct {
	Match []VirtualServiceMatch `yaml:"match"`
	Route []VirtualServiceRoute `yaml:"route"`
}
type VirtualServiceSpec struct {
	Hosts    []string             `yaml:"hosts"`
	Gateways []string             `yaml:"gateways"`
	HTTP     []VirtualServiceHTTP `yaml:"http"`
}
