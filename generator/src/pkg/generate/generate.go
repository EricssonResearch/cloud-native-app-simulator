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
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

const (
	volumeName = "config-data-volume"
	volumePath = "/usr/src/app/config"

	imageName = "app"
	imageURL  = "redis-demo:latest"

	redisImageName = "db"
	redisImageURL  = "k8s.gcr.io/redis:e2e"

	fortioImageName = "fortio"
	fortioImageUrl  = "fortio/fortio"

	workerImageName = "worker"
	workerImageURL  = "redis-demo-worker"

	protocol         = "http"
	redisProtocol    = "tcp"
	fortioProtocol   = "tcp-fortio"
	fortioUiProtocol = "http-fortio-ui"

	defaultExtPort      = 80
	defaultPort         = 5000
	redisDefaultPort    = 6379
	redisTargetPort     = 6379
	fortioDefaultPort   = 8080
	fortioUiDefaultPort = 8089

	uri         = "/"
	fortioUri   = "/echo?"
	fortioUiUri = "/ui"

	replicaNumber   = 1
	numberOfWorkers = 1

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
	Cluster      string  `json:"cluster"`
	Service      string  `json:"service"`
	Endpoint     string  `json:"endpoint"`
	Protocol     string  `json:"protocol"`
	TrafficRatio float32 `json:"traffic_ratio"`
	Requests     string  `json:"requests"`
}
type Endpoints struct {
	Name               string           `json:"name"`
	Protocol           string           `json:"protocol"`
	CpuConsumption     float64          `json:"cpuConsumption"`
	NetworkConsumption float64          `json:"networkConsumption"`
	MemoryConsumption  float64          `json:"memoryConsumption"`
	CalledServices     []CalledServices `json:"calledServices"`
	Requests           string           `json:"requests"`
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
	Node      string      `json:"node"`
	Resources Resources   `json:"resources"`
	Endpoints []Endpoints `json:"endpoints"`
}

type Clusters struct {
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Services  []Services `json:"services"`
}

type Latencies struct {
	Src     string  `json:"src"`
	Dest    string  `json:"dest"`
	Latancy float64 `json:"latancy"`
}

type Config struct {
	Latencies []Latencies `json:"latencies"`
	Clusters  []Clusters  `json:"clusters"`
}

// NumMS the total number of the microservices in the service description file
var NumCluster, NumEP int
var services, clusters, endpoints []string

// Unique return the number of unique elements in the slice of strings
func Unique(strSlice []string) int {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return len(list)
}

// Parse microservice config file, and return a config struct
func Parse(configFilename string) Config {
	configFile, err := os.Open(configFilename)
	configFileByteValue, _ := ioutil.ReadAll(configFile)

	if err != nil {
		fmt.Println(err)
	}

	var loaded_config Config
	json.Unmarshal(configFileByteValue, &loaded_config)
	for i := 0; i < len(loaded_config.Clusters); i++ {
		clusters = append(clusters, loaded_config.Clusters[i].Name)
		NumCluster = NumCluster + 1
		for k := 0; k < len(loaded_config.Clusters[i].Services); k++ {
			services = append(services, loaded_config.Clusters[i].Services[k].Name)
			for l := 0; l < len(loaded_config.Clusters[i].Services[k].Endpoints); l++ {
				endpoints = append(endpoints, loaded_config.Clusters[i].Services[k].Endpoints[l].Name)
				NumEP = NumEP + 1
			}
		}
	}

	fmt.Println("All clusters: ", clusters)
	fmt.Println("Number of clusters: ", NumCluster)
	fmt.Println("---------------")
	fmt.Println("All Services: ", services)
	fmt.Println("Number of services (unique): ", Unique(services))
	fmt.Println("---------------")
	fmt.Println("All endpoints: ", endpoints)
	fmt.Println("Number of endpoints: ", NumEP)

	return loaded_config
}

func Create(config Config, readinessProbe int) {
	path, _ := os.Getwd()
	path = path + "/k8s"
	for i := 0; i < len(config.Clusters); i++ {
		directory := fmt.Sprintf(path+"/%s", config.Clusters[i].Name)
		c_id := fmt.Sprintf("%s", config.Clusters[i].Name)
		os.Mkdir(directory, 0777)
		for k := 0; k < len(config.Clusters[i].Services); k++ {
			serv := config.Clusters[i].Services[k].Name
			resources := Resources(config.Clusters[i].Services[k].Resources)
			nodeAffinity := config.Clusters[i].Services[k].Node

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
			serv_endpoints := []Endpoints(config.Clusters[i].Services[k].Endpoints)
			serv_ep_json, err := json.Marshal(serv_endpoints)
			if err != nil {
				panic(err)
			}
			manifestFilePath := fmt.Sprintf(directory+"/%s.yaml", serv)
			manifests := make([]string, 0, 1)
			appendManifest := func(manifest interface{}) error {
				yamlDoc, err := yaml.Marshal(manifest)
				if err != nil {
					return err
				}
				manifests = append(manifests, string(yamlDoc))
				return nil
			}
			configmap = s.CreateConfig("config-"+serv, "config-"+serv, c_id, config.Clusters[i].Namespace, string(serv_ep_json))
			appendManifest(configmap)
			if nodeAffinity == "" {
				deployment := s.CreateDeployment(serv, serv, c_id, replicaNumber, config.Clusters[i].Name, c_id, config.Clusters[i].Namespace,
					defaultPort, redisDefaultPort, fortioDefaultPort, imageName, imageURL, redisImageName,
					redisImageURL, workerImageName, workerImageURL, fortioImageName, fortioImageUrl, volumePath, volumeName, "config-"+serv, readinessProbe,
					resources.Requests.Cpu, resources.Requests.Memory, resources.Limits.Cpu, resources.Limits.Memory)
				appendManifest(deployment)
			} else {
				deployment := s.CreateDeploymentWithAffinity(serv, serv, c_id, replicaNumber, config.Clusters[i].Name, c_id, config.Clusters[i].Namespace,
					defaultPort, redisDefaultPort, fortioDefaultPort, imageName, imageURL, redisImageName,
					redisImageURL, workerImageName, workerImageURL, fortioImageName, fortioImageUrl, volumePath, volumeName, "config-"+serv, readinessProbe,
					resources.Requests.Cpu, resources.Requests.Memory, resources.Limits.Cpu, resources.Limits.Memory, nodeAffinity)
				appendManifest(deployment)
			}

			service = s.CreateService(serv, serv, protocol, fortioProtocol, fortioUiProtocol, uri, fortioUri, fortioUiUri, c_id, config.Clusters[i].Namespace, defaultExtPort, defaultPort, fortioDefaultPort, fortioUiDefaultPort)
			appendManifest(service)

			yamlDocString := strings.Join(manifests, "---\n")
			ioutil.WriteFile(manifestFilePath, []byte(yamlDocString), 0644)

		}
	}
}
