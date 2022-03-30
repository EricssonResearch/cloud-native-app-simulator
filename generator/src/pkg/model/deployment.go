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

type DeploymentInstance struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Labels    struct {
			App     string `yaml:"app,omitempty"`
			Cluster string `yaml:"version,omitempty"`
		} ` yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Selector struct {
			MatchLabels struct {
				App     string `yaml:"app"`
				Cluster string `yaml:"version"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Replicas int `yaml:"replicas"`
		Template struct {
			Metadata struct {
				Labels struct {
					App     string `yaml:"app"`
					Cluster string `yaml:"version"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec specInstance `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type specInstance struct {
	NodeName       string              `yaml:"nodeName,omitempty"`
	ServiceAccount string              `yaml:"serviceAccountName,omitempty"`
	Containers     []ContainerInstance `yaml:"containers"`
	Volumes        []VolumeInstance    `yaml:"volumes"`
}

type VolumeInstance struct {
	Name      string `yaml:"name"`
	ConfigMap struct {
		Name string `yaml:"name"`
	} `yaml:"configMap"`
}
type ContainerInstance struct {
	Name            string                    `yaml:"name"`
	Image           string                    `yaml:"image"`
	ImagePullPolicy string                    `yaml:"imagePullPolicy"`
	Env             []EnvInstance             `yaml:"env"`
	Ports           []ContainerPortInstance   `yaml:"ports"`
	Volumes         []ContainerVolumeInstance `yaml:"volumeMounts"`
	ReadinessProbe  ReadinessProbeInstance    `yaml:"readinessProbe,omitempty"`
	Resources       ResourcesInstance         `yaml:"resources"`
}

type ContainerPortInstance struct {
	ContainerPort int `yaml:"containerPort"`
}

type EnvInstance struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ContainerVolumeInstance struct {
	MountPath string `yaml:"mountPath,omitempty"`
	MountName string `yaml:"name,omitempty"`
}

type ReadinessProbeInstance struct {
	HttpGet struct {
		Path string `yaml:"path,omitempty"`
		Port int    `yaml:"port,omitempty"`
	} `yaml:"httpGet,omitempty"`
	Exec struct {
		Command []string `yaml:"command,flow,omitempty"`
	} `yaml:"exec,omitempty"`
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`
	PeriodSeconds       int `yaml:"periodSeconds"`
}
type ResourcesInstance struct {
	ResourceLimits struct {
		Cpu    string `yaml:"cpu"`
		Memory string `yaml:"memory"`
	} `yaml:"limits"`
	ResourceRequests struct {
		Cpu    string `yaml:"cpu"`
		Memory string `yaml:"memory"`
	} `yaml:"requests"`
}
