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

type ConfigMapInstance struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			Name    string `yaml:"name"`
			Cluster string `yaml:"version,omitempty"`
		} `yaml:"labels"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Data struct {
		Config  string `yaml:"conf.json"`
		Service string `yaml:"service.proto"`
	} `yaml:"data"`
}

type ConfigMap struct {
	Processes int        `json:"processes"`
	Threads   int        `json:"threads"`
	Logging   bool       `json:"logging"`
	Endpoints []Endpoint `json:"endpoints"`
}
