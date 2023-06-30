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
	"cloud-native-app-simulator/emulator/src/util"
	"cloud-native-app-simulator/model"

	"cloud-native-app-simulator/emulator/src/restful"

	"fmt"
	"sync"
)

func main() {
	configMap, err := util.LoadConfigMap()
	if err != nil {
		fmt.Println("Using default config map")
		configMap = util.DefaultConfigMap()
	}

	processes := util.SetMaxProcesses(configMap)
	fmt.Printf("Application emulator started at :5000, %s\n\n", processes)

	wg := sync.WaitGroup{}

	// TODO: Check if protocol is HTTP
	httpEndpoints := make(chan model.Endpoint)
	go restful.HTTP(httpEndpoints, &wg)
	wg.Add(1)

	fmt.Println("Endpoints:")
	for _, endpoint := range configMap.Endpoints {
		// Only HTTP is supported right now
		if endpoint.Protocol == "http" {
			fmt.Println("*", endpoint.Protocol, endpoint.Name)
			httpEndpoints <- endpoint
		} else {
			fmt.Println("x", endpoint.Protocol, endpoint.Name)
		}
	}

	close(httpEndpoints)
	wg.Wait()
}
