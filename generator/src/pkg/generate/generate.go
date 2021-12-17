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
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	volumeName = "config-data-volume"
	volumePath = "/usr/src/app/config"

	imageName = "app"
	imageURL  = "redis-demo:latest"

	redisImageName = "db"
	redisImageURL  = "k8s.gcr.io/redis:e2e"

	fortioImageName = "fortio"
	fortioImageUrl = "fortio/fortio"

	workerImageName = "worker"
	workerImageURL  = "redis-demo-worker"

	protocol      = "http"
	redisProtocol = "tcp"
	fortioProtocol = "tcp-fortio"
	fortioUiProtocol = "http-fortio-ui"

	namespace = "edge-namespace"

	defaultExtPort = 80
	defaultPort = 5000
	redisDefaultPort = 6379
	redisTargetPort  = 6379
	fortioDefaultPort = 8080
	fortioUiDefaultPort = 8089

	uri = "/"
	fortioUri = "/echo?"
	fortioUiUri = "/ui"

	replicaNumber   = 1
	numberOfWorkers = 1
)

var (
	config           model.ConfigMapInstance
	deployment       model.DeploymentInstance
	service          model.ServiceInstance
	serviceAccount   model.ServiceAccountInstance
	virtualService   model.VirtualServiceInstance
	workerDeployment model.DeploymentInstance
)

type ConfigData struct {
	Hop      ChainTypes `json:"Hop"`
	Hostname string     `json:"Hostname"`
}

type ChainTypes struct {
	Chain1 []string `json:"1"`
	Chain2 []string `json:"2"`
	Chain3 []string `json:"3"`
	Chain4 []string `json:"4"`
	Chain5 []string `json:"5"`
	Chain6 []string `json:"6"`
	Chain7 []string `json:"7"`
	Chain8 []string `json:"8"`
}

type ServiceChains struct {
	ServiceID string   `json:"serviceId"`
	Chains    []Chains `json:"chains"`
	Number    int      `json:"number"`
}

type Chains struct {
	ChainID  string   `json:"chainid"`
	Latency  float64  `json:"latency"`
	Services []string `json:"microservices"`
}
type Placement struct {
    ServiceID   string  `json:"service_id"`
    FinalPlacement []struct {
        ClusterID string `json:"cluster_id"`
        Service   []string `json:"service_list"`
	} `json:"final_placement"`
}

// NumMS the total number of the microservices in the service description file
var NumMS int
var mlist []string

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

// Parse microservice chain and cluster file, and return a mapped placement struct
func Parse(chainFilename, clusterFilename string) (m map[string]map[string][]string, p Placement) {

	// filePath := fmt.Sprintf("chains/%s", chainFilename)

	chainFile, err := os.Open(chainFilename)

	if err != nil {
		fmt.Println(err)
	}

	chainFileByteValue, _ := ioutil.ReadAll(chainFile)

	// clusterFilePath := fmt.Sprintf("clusters/%s", clusterFilename)

	clusterFile, e := os.Open(clusterFilename)

	if e != nil {
		fmt.Println(e)
	}

	clusterFileByteValue, _ := ioutil.ReadAll(clusterFile)

	var serviceChains ServiceChains
	var placement Placement

	json.Unmarshal(chainFileByteValue, &serviceChains)
	json.Unmarshal(clusterFileByteValue, &placement)

	m = map[string]map[string][]string{}

	for i := 0; i < len(serviceChains.Chains); i++ {
		w := map[string][]string{}
		for k := 1; k < len(serviceChains.Chains[i].Services); k++ {
			data := serviceChains.Chains[i].Services[k]
			w[serviceChains.Chains[i].Services[k-1]] = append(w[serviceChains.Chains[i].Services[k-1]], data)
		}
		m[serviceChains.Chains[i].ChainID] = w
		mlist = append(mlist, serviceChains.Chains[i].Services[:]...)
	}
	fmt.Println("all ms: ", mlist)
	NumMS = Unique(mlist)
	fmt.Println("the unique number of ms is: ", NumMS)

	return m, placement

}
func Create(m map[string]map[string][]string, placement Placement, readinessProbe int) {

	path, _ := os.Getwd()
	path = path + "/k8s"
	for i := 0; i < len(placement.FinalPlacement); i++ {
		directory := fmt.Sprintf(path+"/%s", placement.FinalPlacement[i].ClusterID)
		c_id := fmt.Sprintf("%s", placement.FinalPlacement[i].ClusterID)
		os.Mkdir(directory, 0777)
		for j := 0; j < len(placement.FinalPlacement[i].Service); j++ {
			configData := ConfigData{}
			serv := placement.FinalPlacement[i].Service[j]
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
			configData = ConfigData{Hop: ChainTypes{Chain1: m["1"][serv], Chain2: m["2"][serv], Chain3: m["3"][serv], Chain4: m["4"][serv], Chain5: m["5"][serv],
				Chain6: m["6"][serv], Chain7: m["7"][serv], Chain8: m["8"][serv]}, Hostname: serv}
			d, _ := json.Marshal(configData)
			data := fmt.Sprintf("%s", string(d))

			config = s.CreateConfig("config-"+serv, "config-"+serv, c_id, namespace, data)
			appendManifest(config)
			deployment = s.CreateDeployment(serv, serv, c_id, replicaNumber, serv, c_id, namespace,
			    defaultPort, redisDefaultPort, fortioDefaultPort, imageName, imageURL, redisImageName,
				redisImageURL, workerImageName, workerImageURL, fortioImageName, fortioImageUrl, volumePath, volumeName, "config-"+serv, readinessProbe)
			appendManifest(deployment)
			service = s.CreateService(serv, serv, protocol, fortioProtocol, fortioUiProtocol, uri, fortioUri, fortioUiUri, c_id, namespace, defaultExtPort, defaultPort, fortioDefaultPort, fortioUiDefaultPort)
			appendManifest(service)

			yamlDocString := strings.Join(manifests, "---\n")
			ioutil.WriteFile(manifestFilePath, []byte(yamlDocString), 0644)
			}
	}

}
