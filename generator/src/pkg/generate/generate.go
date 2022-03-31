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
package generate

import (
	"application-generator/src/pkg/model"
	s "application-generator/src/pkg/service"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const (
	volumeName = "config-data-volume"
	volumePath = "/usr/src/app/config"

	imageName = "app"
	imageURL  = "app-demo:latest"

	defaultExtPort = 80
	defaultPort    = 5000

	uri = "/"

	replicaNumber = 1

	requestsCPUDefault    = "500m"
	requestsMemoryDefault = "256M"
	limitsCPUDefault      = "1000m"
	limitsMemoryDefault   = "1024M"
)

var (
	configmap        model.ConfigMapInstance
	deployment       model.DeploymentInstance
	service          model.ServiceInstance
	serviceAccount   model.ServiceAccountInstance
	virtualService   model.VirtualServiceInstance
	workerDeployment model.DeploymentInstance
)

type CalledServices struct {
	Service             string  `json:"service"`
	Port                string  `json:"port"`
	Endpoint            string  `json:"endpoint"`
	Protocol            string  `json:"protocol"`
	TrafficForwardRatio float32 `json:"traffic_forward_ratio"`
}
type Endpoints struct {
	Name               string           `json:"name"`
	Protocol           string           `json:"protocol"`
	CpuConsumption     float64          `json:"cpu_consumption"`
	NetworkConsumption float64          `json:"network_consumption"`
	MemoryConsumption  float64          `json:"memory_consumption"`
	ForwardRequests    string           `json:"forward_requests"`
	CalledServices     []CalledServices `json:"called_services"`
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
type Services struct {
	Name      string      `json:"name"`
	Clusters  []Clusters  `json:"clusters"`
	Resources Resources   `json:"resources"`
	Processes int         `json:"processes"`
	Endpoints []Endpoints `json:"endpoints"`
}

type Clusters struct {
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	Node      string `json:"node,omitempty"`
}

type Latencies struct {
	Src     string  `json:"src"`
	Dest    string  `json:"dest"`
	Latancy float64 `json:"latency"`
}

type Config struct {
	Latencies []Latencies `json:"cluster_latencies"`
	Services  []Services  `json:"services"`
}

type ConfigMap struct {
	Processes int         `json:"processes"`
	Endpoints []Endpoints `json:"endpoints"`
}

// the slices to store services, cluster and endpoints for counting and printing
var services, clusters, endpoints []string

// Unique return unique elements in the slice of strings
func Unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Parse microservice config file, and return a config struct
func Parse(configFilename string) (Config, []string) {
	configFile, err := os.Open(configFilename)
	configFileByteValue, _ := ioutil.ReadAll(configFile)

	if err != nil {
		fmt.Println(err)
	}

	var loaded_config Config
	json.Unmarshal(configFileByteValue, &loaded_config)
	for i := 0; i < len(loaded_config.Services); i++ {
		services = append(services, loaded_config.Services[i].Name)
		for j := 0; j < len(loaded_config.Services[i].Clusters); j++ {
			clusters = append(clusters, loaded_config.Services[i].Clusters[j].Cluster)
		}
		for k := 0; k < len(loaded_config.Services[i].Endpoints); k++ {
			endpoints = append(endpoints, loaded_config.Services[i].Endpoints[k].Name)
		}

	}

	fmt.Println("All clusters: ", Unique(clusters))
	fmt.Println("Number of clusters: ", len(Unique(clusters)))
	fmt.Println("---------------")
	fmt.Println("All Services: ", Unique(services))
	fmt.Println("Number of services (unique): ", len(Unique(services)))
	fmt.Println("---------------")
	fmt.Println("All endpoints: ", Unique(endpoints))
	fmt.Println("Number of endpoints: ", len(Unique(endpoints)))
	return loaded_config, clusters
}

func Create(config Config, readinessProbe int, clusters []string) {
	path, _ := os.Getwd()
	proto_temp, _ := template.ParseFiles(path + "/template/service.tmpl")
	path = path + "/k8s"

	for i := 0; i < len(clusters); i++ {
		directory := fmt.Sprintf(path+"/%s", clusters[i])
		os.Mkdir(directory, 0777)
	}
	var proto_temp_filled_byte bytes.Buffer
	err := proto_temp.Execute(&proto_temp_filled_byte, config.Services)
	if err != nil {
		panic(err)
	}
	proto_temp_filled := proto_temp_filled_byte.String()
	for i := 0; i < len(config.Services); i++ {
		serv := config.Services[i].Name
		resources := Resources(config.Services[i].Resources)
		protocol := config.Services[i].Endpoints[0].Protocol

		if resources.Limits.Cpu == "" {
			resources.Limits.Cpu = limitsCPUDefault
		}
		if resources.Limits.Memory == "" {
			resources.Limits.Memory = limitsMemoryDefault
		}
		if resources.Requests.Cpu == "" {
			resources.Requests.Cpu = requestsCPUDefault
		}
		if resources.Requests.Memory == "" {
			resources.Requests.Memory = requestsMemoryDefault
		}

		cm_data := &ConfigMap{
			Processes: int(config.Services[i].Processes),
			Endpoints: []Endpoints(config.Services[i].Endpoints),
		}

		serv_json, err := json.Marshal(cm_data)
		if err != nil {
			panic(err)
		}

		for j := 0; j < len(config.Services[i].Clusters); j++ {
			directory := config.Services[i].Clusters[j].Cluster
			directory_path := fmt.Sprintf(path+"/%s", directory)
			c_id := fmt.Sprintf("%s", config.Services[i].Clusters[j].Cluster)
			nodeAffinity := config.Services[i].Clusters[j].Node
			namespace := config.Services[i].Clusters[j].Namespace
			manifestFilePath := fmt.Sprintf(directory_path+"/%s.yaml", serv)
			manifests := make([]string, 0, 1)
			appendManifest := func(manifest interface{}) error {
				yamlDoc, err := yaml.Marshal(manifest)
				if err != nil {
					return err
				}
				manifests = append(manifests, string(yamlDoc))
				return nil
			}
			configmap = s.CreateConfig("config-"+serv, "config-"+serv, c_id, namespace, string(serv_json), proto_temp_filled)
			appendManifest(configmap)

			deployment := s.CreateDeployment(serv, serv, c_id, replicaNumber, serv, c_id, namespace,
				defaultPort, imageName, imageURL, volumePath, volumeName, "config-"+serv, readinessProbe,
				resources.Requests.Cpu, resources.Requests.Memory, resources.Limits.Cpu, resources.Limits.Memory,
				nodeAffinity, protocol)
			appendManifest(deployment)

			service = s.CreateService(serv, serv, protocol, uri, c_id, namespace, defaultExtPort, defaultPort)
			appendManifest(service)

			yamlDocString := strings.Join(manifests, "---\n")
			err := ioutil.WriteFile(manifestFilePath, []byte(yamlDocString), 0644)
			if err != nil {
				fmt.Print(err)
				return
			}

		}
	}
}
