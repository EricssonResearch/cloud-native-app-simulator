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

type RedisDeploymentInstance struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			App string `yaml:"app"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Selector struct {
			MatchLabels struct {
				App  string `yaml:"app"`
				Role string `yaml:"role"`
				Tier string `yaml:"tier"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Replicas int `yaml:"replicas"`
		Template struct {
			Metadata struct {
				Labels struct {
					App  string `yaml:"app"`
					Role string `yaml:"role"`
					Tier string `yaml:"tier"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []RedisContainerInstance `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type RedisContainerInstance struct {
	Name            string                  `yaml:"name"`
	Image           string                  `yaml:"image"`
	ImagePullPolicy string                  `yaml:"imagePullPolicy"`
	Ports           []ContainerPortInstance `yaml:"ports"`
	Resources       struct {
		Requests struct {
			CPU    string `yaml:"cpu"`
			Memory string `yaml:"memory"`
		}
	}
}

type RedisServiceInstance struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			App  string `yaml:"app"`
			Role string `yaml:"role"`
			Tier string `yaml:"tier"`
		}
	} `yaml:"metadata"`
	Spec struct {
		Selector struct {
			App  string `yaml:"app"`
			Role string `yaml:"role"`
			Tier string `yaml:"tier"`
		} `yaml:"selector"`
		Ports []ServicePortInstance `yaml:"ports"`
	} `yaml:"spec"`
}
