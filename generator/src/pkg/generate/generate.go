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
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

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

	loaded_config := s.CreateFileConfig()

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
		resources := s.CreateInputResources()
		resources = config.Services[i].Resources
		protocol := config.Services[i].Endpoints[0].Protocol

		if resources.Limits.Cpu == "" {
			resources.Limits.Cpu = s.LimitsCPUDefault
		}
		if resources.Limits.Memory == "" {
			resources.Limits.Memory = s.LimitsMemoryDefault
		}
		if resources.Requests.Cpu == "" {
			resources.Requests.Cpu = s.RequestsCPUDefault
		}
		if resources.Requests.Memory == "" {
			resources.Requests.Memory = s.RequestsMemoryDefault
		}

		readinessProbe := config.Services[i].ReadinessProbe
		if readinessProbe == 0 {
			readinessProbe = s.SvcReadinessProbeDefault
		}

		processes := config.Services[i].Processes
		if processes == 0 {
			processes = s.SvcProcessesDefault
		}

		threads := config.Services[i].Threads
		if threads == 0 {
			processes = s.SvcThreadsDefault
		}

		cm_data := s.CreateConfigMap(processes, threads, config.Services[i].Endpoints)

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
			configmap := s.CreateConfig("config-"+serv, "config-"+serv, c_id, namespace, string(serv_json), proto_temp_filled)
			appendManifest(configmap)

			deployment := s.CreateDeployment(serv, serv, c_id, s.ReplicaNumber, serv, c_id, namespace,
				s.DefaultPort, s.ImageName, s.ImageURL, s.VolumePath, s.VolumeName, "config-"+serv, readinessProbe,
				resources.Requests.Cpu, resources.Requests.Memory, resources.Limits.Cpu, resources.Limits.Memory,
				nodeAffinity, protocol)
			appendManifest(deployment)

			service := s.CreateService(serv, serv, protocol, s.Uri, c_id, namespace, s.DefaultExtPort, s.DefaultPort)
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

func CreateJsonInput(userConfig model.UserConfig) string {
	path, _ := os.Getwd()
	path = path + "/input/" + userConfig.OutputFileName

	rand.Seed(time.Now().UnixNano())

	inputConfig := s.CreateFileConfig()

	// TODO: Generate cluster latencies

	// Generating random services
	serviceNumber := rand.Intn(userConfig.SvcMaxNumber) + 1
	for i := 1; i <= serviceNumber; i++ {
		service := s.CreateInputService()

		service.Name = s.SvcNamePrefix + strconv.Itoa(i)

		// Randomly associating services to clusters
		replicaNumber := rand.Intn(userConfig.SvcReplicaMaxNumber) + 1
		for j := 1; j <= replicaNumber; j++ {
			cluster := s.CreateInputCluster()

			cRIndex := rand.Intn(len(userConfig.Clusters))
			cluster.Cluster = userConfig.Clusters[cRIndex]

			nRIndex := rand.Intn(len(userConfig.Namespaces))
			cluster.Namespace = userConfig.Namespaces[nRIndex]

			service.Clusters = append(service.Clusters, cluster)
		}

		resources := s.CreateInputResources()
		service.Resources = resources

		service.Processes = s.SvcProcessesDefault
		service.Threads = s.SvcThreadsDefault
		service.ReadinessProbe = s.SvcReadinessProbeDefault

		// Randomly generating service endpoints
		endpointNumber := rand.Intn(userConfig.SvcEpMaxNumber) + 1
		for k := 1; k <= endpointNumber; k++ {
			endpoint := s.CreateInputEndpoint()

			endpoint.Name = s.EpNamePrefix + strconv.Itoa(k)

			// Randomly generating called services
			// NOTE: Last service does not call any service to ensure the sequence of calls ends
			if i < serviceNumber {
				// NOTE: Services only call subsequent services to avoid endless loops
				calledServiceNumber := rand.Intn(serviceNumber-i+1) + i // (max - min + 1) + min
				for n := i + 1; n <= calledServiceNumber; n++ {
					calledService := s.CreateInputCalledSvc()

					calledService.Service = s.SvcNamePrefix + strconv.Itoa(n)
					// NOTE: Always calling the first endpoint of the called service
					calledService.Endpoint = s.EpNamePrefix + "1"

					endpoint.CalledServices = append(endpoint.CalledServices, calledService)
				}
			}

			service.Endpoints = append(service.Endpoints, endpoint)
		}

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
