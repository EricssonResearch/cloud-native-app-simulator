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
	Endpoints []model.Endpoint	`json:"endpoints"`
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
func Parse(configFilename string) (model.FileConfig, []string) {
	configFile, err := os.Open(configFilename)
	configFileByteValue, _ := ioutil.ReadAll(configFile)

	if err != nil {
		fmt.Println(err)
	}

	var loaded_config model.FileConfig
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

func CreateK8sYaml(config model.FileConfig, clusters []string) {
	path, _ := os.Getwd()
	path = path + "/k8s"

	for i := 0; i < len(clusters); i++ {
		directory := fmt.Sprintf(path+"/%s", clusters[i])
		os.Mkdir(directory, 0777)
	}

	for i := 0; i < len(config.Services); i++ {
		serv := config.Services[i].Name
		resources := model.Resources(config.Services[i].Resources)

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
		if readinessProbe == 0 {
			readinessProbe = serviceReadinessProbeDefault
		}

		processes := config.Services[i].Processes
		if processes == 0 {
			processes = serviceProcessesDefault
		}

		cm_data := &ConfigMap{
			Processes: processes,
			Endpoints: []model.Endpoint(config.Services[i].Endpoints),
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

func CreateJsonInput(clusterConfig model.ClusterConfig) (string) {
	path, _ := os.Getwd()
	path = path + "/input/new_description_test.json"

	var inputConfig model.FileConfig

	// TODO: Generate cluster latencies

	// Generating random services
	// TODO: Replace this hard-coded number of services by the ones given by the user
	serviceNumber := rand.Intn(10)
	for i := 1; i <= serviceNumber; i++ {
		var service model.Service

		service.Name = "service" + strconv.Itoa(i)

		// Randomly associating services to clusters
		// TODO: Replace this hard-coded number of service replicas by the ones given by the user
		replicaNumber := rand.Intn(5)
		for j := 1; j <= replicaNumber; j++ {
			var cluster model.Cluster

			rIndex := rand.Intn(len(clusterConfig.Clusters))
			cluster.Cluster = clusterConfig.Clusters[rIndex]

			rIndex = rand.Intn(len(clusterConfig.Namespaces))
			cluster.Namespace = clusterConfig.Namespaces[rIndex]

			service.Clusters = append(service.Clusters, cluster)
		}

		var resources model.Resources
		resources.ResourceLimits.Cpu = limitsCPUDefault
		resources.ResourceLimits.Memory = limitsMemoryDefault
		resources.ResourceRequests.Cpu = requestsCPUDefault
		resources.ResourceRequests.memory = requestsMemoryDefault
		service.Resources = resources

		service.Processes = serviceProcessesDefault
		service.ReadinessProbe = serviceReadinessProbeDefault

		inputConfig.Services = append(inputConfig.Services, service)
	}

	input_json, err := json.MarshalIndent(inputConfig, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, input_json, 0644)
	if err != nil {
		fmt.Print(err)
	}

	return path
}