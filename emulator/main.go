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
	"cloud-native-app-simulator/model"

	"encoding/json"
	"io/ioutil"
	"os"

	"fmt"
	"net/http"
)

const httpPort = 5000

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

	fmt.Println("Read config:", configFilename)
	fmt.Println("Endpoints:", len(configMap.Endpoints))

	// Placeholder, start a HTTP server at :5000
	fmt.Println("Application emulator started at :5000")

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprint(writer, "{\"status\": \"ok\"}\n")
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)

	if err != nil {
		panic(err)
	}
}
