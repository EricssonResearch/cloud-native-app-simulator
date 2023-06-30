/*
Copyright 2023 Telefonaktiebolaget LM Ericsson AB

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

package main

import (
	"cloud-native-app-simulator/emulator/src/restful"
	"cloud-native-app-simulator/model"
	"sync"

	"runtime"

	"encoding/json"
	"io/ioutil"
	"os"

	"fmt"
)

func loadConfigMap(filename string) (*model.ConfigMap, error) {
	configFile, err := os.Open(filename)
	configFileByteValue, _ := ioutil.ReadAll(configFile)

	if err != nil {
		return nil, err
	}

	inputConfig := &model.ConfigMap{}
	err = json.Unmarshal(configFileByteValue, inputConfig)

	if err != nil {
		return nil, err
	}

	return inputConfig, nil
}

func main() {
	configFilename := os.Getenv("CONF")
	configMap, err := loadConfigMap(configFilename)

	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(configMap.Processes)

	processes := fmt.Sprintf("%d processes", configMap.Processes)
	if configMap.Processes == 1 {
		processes = fmt.Sprintf("%d process", configMap.Processes)
	}

	serverWaitGroup := sync.WaitGroup{}
	httpEndpointChannel := make(chan model.Endpoint)

	go restful.HTTP(httpEndpointChannel, &serverWaitGroup)
	serverWaitGroup.Add(1)

	fmt.Println("Application emulator started at :5000,", processes)
	fmt.Println()

	fmt.Println("Endpoints:")

	for _, endpoint := range configMap.Endpoints {
		// Only HTTP is supported right now
		if endpoint.Protocol == "http" {
			fmt.Println("*", endpoint.Protocol, endpoint.Name)
			httpEndpointChannel <- endpoint
		} else {
			fmt.Println("x", endpoint.Protocol, endpoint.Name)
		}
	}

	close(httpEndpointChannel)
	serverWaitGroup.Wait()
}
