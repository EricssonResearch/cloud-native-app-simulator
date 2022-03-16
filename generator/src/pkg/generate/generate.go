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
	"math/rand"
)

const (
	volumeName = "config-data-volume"
	volumePath = "/usr/src/app/config"

	imageName = "app"
	imageURL  = "app-demo:latest"

	protocol = "http"

	defaultExtPort = 80
	defaultPort    = 5000

	uri = "/"

	replicaNumber = 1

	requestsCPUDefault    = "500m"
	requestsMemoryDefault = "256M"
	limitsCPUDefault      = "1000m"
	limitsMemoryDefault   = "1024M"

	serviceProcessesDefault = 2
	serviceReadinessProbeDefault = 5
)

var (
	configmap        model.ConfigMapInstance
	deployment       model.DeploymentInstance
	service          model.ServiceInstance
	serviceAccount   model.ServiceAccountInstance
	virtualService   model.VirtualServiceInstance
	workerDeployment model.DeploymentInstance
)

type ConfigMap struct {
	Processes	int					`json:"processes"`
	Endpoints []Endpoints	`json:"endpoints"`
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
func Parse(configFilename string) (FileConfig, []string) {
	configFile, err := os.Open(configFilename)
	configFileByteValue, _ := ioutil.ReadAll(configFile)

	if err != nil {
		fmt.Println(err)
	}

	var loaded_config FileConfig
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

func CreateK8sYaml(config FileConfig, clusters []string) {
	path, _ := os.Getwd()
	path = path + "/k8s"

	for i := 0; i < len(clusters); i++ {
		directory := fmt.Sprintf(path+"/%s", clusters[i])
		os.Mkdir(directory, 0777)
	}

	for i := 0; i < len(config.Services); i++ {
		serv := config.Services[i].Name
		resources := Resources(config.Services[i].Resources)

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

		readinessProbe := config.Services[i].ReadinessProbe
		if readinessProbe == NULL {
			readinessProbe = serviceReadinessProbeDefault
		}

		processes := config.Services[i].Processes
		if processes == NULL {
			processes = serviceProcessesDefault
		}

		cm_data := &ConfigMap{
			Processes: processes,
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
			configmap = s.CreateConfig("config-"+serv, "config-"+serv, c_id, namespace, string(serv_json))
			appendManifest(configmap)
			if nodeAffinity == "" {
				deployment := s.CreateDeployment(serv, serv, c_id, replicaNumber, serv, c_id, namespace,
					defaultPort, imageName, imageURL, volumePath, volumeName, "config-"+serv, readinessProbe,
					resources.Requests.Cpu, resources.Requests.Memory, resources.Limits.Cpu, resources.Limits.Memory)

				appendManifest(deployment)
			} else {
				deployment := s.CreateDeploymentWithAffinity(serv, serv, c_id, replicaNumber, serv, c_id, namespace,
					defaultPort, imageName, imageURL, volumePath, volumeName, "config-"+serv, readinessProbe,
					resources.Requests.Cpu, resources.Requests.Memory, resources.Limits.Cpu, resources.Limits.Memory, nodeAffinity)
				appendManifest(deployment)
			}

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

func CreateJsonInput(clusterConfig ClusterConfig) (string) {
	path, _ := os.Getwd()
	path = path + "/input/new_description_test.json"

	var inputConfig FileConfig

	// TODO: Generate cluster latencies

	// Generating random services
	serviceNumber := rand.Intn(10)
	for i := 1; i <= serviceNumber; i++ {
		var service Service

		service.name = "service" + i

		// Randomly associating services to clusters
		replicaNumber := rand.Intn(5)
		for j := 1; j <= replicaNumber; j++ {
			var cluster Cluster

			rIndex := rand.Intn(len(clusterConfig.clusters))
			cluster.cluster = clusterConfig.clusters[rIndex]

			rIndex = rand.Intn(len(clusterConfig.namespaces))
			cluster.namespace = clusterConfig.namespaces[rIndex]

			service.clusters = append(service.clusters, cluster)
		}

		var resource Resource
		resource.limits.cpu = limitsCPUDefault
		resource.limits.memory = limitsMemoryDefault
		resource.requests.cpu = requestsCPUDefault
		resource.requests.memory = requestsMemoryDefault
		service.resource = resource

		service.processes = serviceProcessesDefault
		service.readinessProbe = serviceReadinessProbeDefault

		inputConfig.services = append(inputConfig.services, service)
	}

	input_json, err := json.MarshalIndent(inputConfig, "", " ")
	if err != nil {
		panic(err)
	}

	err := ioutil.WriteFile(path, input_json, 0644)
	if err != nil {
		fmt.Print(err)
		return
	}

	return path
}