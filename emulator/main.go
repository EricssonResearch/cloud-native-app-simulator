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
	"application-emulator/src/generated"
	"application-emulator/src/server"
	"application-emulator/src/util"
	"log"
	"os"
	"runtime"
)

// Randomly generated string in Dockerfile which is used to make sure the binary is up to date with the configuration
var buildID string

func main() {
	configMap, err := util.LoadConfigMap()
	if err != nil {
		configMap = util.DefaultConfigMap()
	}

	if configMap.BuildID != "" && configMap.BuildID != buildID {
		log.Printf("Build ID mismatch: %s != %s, have you deployed the latest Docker image?", configMap.BuildID, buildID)
	}

	runtime.GOMAXPROCS(configMap.Processes)

	util.LoggingEnabled = configMap.Logging
	if name, ok := os.LookupEnv("SERVICE_NAME"); ok {
		util.ServiceName = name
	}
	util.LogConfiguration(configMap)

	util.GRPCCallGeneratedEndpoint = generated.CallGeneratedEndpoint
	util.GRPCRegisterGeneratedService = generated.RegisterGeneratedService

	if configMap.Protocol == "http" {
		server.HTTP(configMap.Endpoints)
	} else if configMap.Protocol == "grpc" {
		server.GRPC(configMap.Endpoints)
	}
}
