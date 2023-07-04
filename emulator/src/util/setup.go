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

package util

import (
	model "application-model"

	"fmt"
	"io"
	"os"
	"runtime"

	"encoding/json"
)

// For local development, will be removed later
func DefaultConfigMap() *model.ConfigMap {
	return &model.ConfigMap{
		Processes: 4,
		Threads:   4,
		Logging:   false,
		Endpoints: []model.Endpoint{
			{
				Name:              "test-endpoint",
				Protocol:          "http",
				ExecutionMode:     "sequential",
				CpuComplexity:     nil,
				NetworkComplexity: nil,
			},
		},
	}
}

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

// Configure the Go runtime to use the number of processes specified in the config map
func SetMaxProcesses(configMap *model.ConfigMap) string {
	runtime.GOMAXPROCS(configMap.Processes)

	if configMap.Processes == 1 {
		return "1 process"
	} else {
		return fmt.Sprintf("%d processes", configMap.Processes)
	}
}
