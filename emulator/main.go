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
	"application-emulator/src/server"
	"application-emulator/src/util"
	model "application-model"
	"encoding/json"
	"io"
	"os"
	"runtime"
)

// Load the config map from the CONF environment variable
func LoadConfigMap() (*model.ConfigMap, error) {
	configFilename := os.Getenv("CONF")
	configFile, err := os.Open(configFilename)
	configFileByteValue, _ := io.ReadAll(configFile)

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
	configMap, err := LoadConfigMap()
	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(configMap.Processes)

	util.LoggingEnabled = configMap.Logging
	if name, ok := os.LookupEnv("SERVICE_NAME"); ok {
		util.ServiceName = name
	}
	util.LogConfiguration(configMap)

	if configMap.Protocol == "http" {
		server.HTTP(configMap.Endpoints)
	} else if configMap.Protocol == "grpc" {
		server.GRPC(configMap.Endpoints)
	}
}
