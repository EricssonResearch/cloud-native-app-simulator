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

type CalledService struct {
	Service             string `json:"service"`
	Port                string `json:"port"`
	Endpoint            string `json:"endpoint"`
	Protocol            string `json:"protocol"`
	TrafficForwardRatio int    `json:"traffic_forward_ratio"`
	RequestPayloadSize  int    `json:"request_payload_size"`
}

type CpuComplexity struct {
	ExecutionTime string `json:"execution_time"`
	Method        string `json:"method"`
	Workers       int    `json:"workers"`
	CpuAffinity   []int  `json:"cpu_affinity"`
	CpuLoad       string `json:"cpu_load"`
}

type MemoryComplexity struct {
	ExecutionTime string `json:"execution_time"`
	Method        string `json:"method"`
	Workers       int    `json:"workers"`
	BytesLoad     string `json:"bytes_load"`
}

type NetworkComplexity struct {
	ForwardRequests     string          `json:"forward_requests"`
	ResponsePayloadSize int             `json:"response_payload_size"`
	CalledServices      []CalledService `json:"called_services"`
}

type Endpoint struct {
	Name              string            `json:"name"`
	Protocol          string            `json:"protocol"`
	ExecutionMode     string            `json:"execution_mode"`
	CpuComplexity     *CpuComplexity    `json:"cpu_complexity,omitempty"`
	MemoryComplexity  *MemoryComplexity `json:"memory_complexity,omitempty"`
	NetworkComplexity NetworkComplexity `json:"network_complexity"`
}

type ResourceLimits struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type ResourceRequests struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Resources struct {
	Limits   ResourceLimits   `json:"limits"`
	Requests ResourceRequests `json:"requests"`
}

type Service struct {
	Name           string     `json:"name"`
	Clusters       []Cluster  `json:"clusters"`
	Resources      Resources  `json:"resources"`
	Processes      int        `json:"processes"`
	Threads        int        `json:"threads"`
	ReadinessProbe int        `json:"readiness_probe"`
	Endpoints      []Endpoint `json:"endpoints"`
}

type Cluster struct {
	Cluster   string `json:"cluster"`
	Replicas  int    `json:"replicas,omitempty"`
	Namespace string `json:"namespace"`
	Node      string `json:"node,omitempty"`
}

type ClusterLatency struct {
	Src     string  `json:"src"`
	Dest    string  `json:"dest"`
	Latency float64 `json:"latency"`
}

type Setting struct {
	Logging bool `json:"logging"`
}

type FileConfig struct {
	ClusterLatencies []ClusterLatency `json:"cluster_latencies"`
	Services         []Service        `json:"services"`
	Settings         Setting          `json:"settings,omitempty"`
}

type UserConfig struct {
	Clusters            []string
	Namespaces          []string
	SvcMaxNumber        int
	SvcReplicaMaxNumber int
	SvcEpMaxNumber      int
	OutputFileName      string
}
