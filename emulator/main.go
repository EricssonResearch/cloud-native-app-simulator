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
	"os"
	"runtime"
	"sync"
)

func main() {
	configMap, err := util.LoadConfigMap()
	if err != nil {
		configMap = util.DefaultConfigMap()
	}

	runtime.GOMAXPROCS(configMap.Processes)

	util.LoggingEnabled = configMap.Logging
	if name, ok := os.LookupEnv("SERVICE_NAME"); ok {
		util.ServiceName = name
	}
	util.LogConfiguration(configMap)

	wg := sync.WaitGroup{}
	httpEndpoints := make(chan model.Endpoint)
	grpcEndpoints := make(chan model.Endpoint)
	grpcStarted := false

	// Always launch HTTP server (required for readiness probe)
	go func() {
		defer wg.Done()
		server.HTTP(httpEndpoints)
	}()
	wg.Add(1)

	for _, endpoint := range configMap.Endpoints {
		if endpoint.Protocol == "http" {
			httpEndpoints <- endpoint
		} else if endpoint.Protocol == "grpc" {
			// Launch gRPC server on demand
			if !grpcStarted {
				go func() {
					defer wg.Done()
					server.GRPC(grpcEndpoints)
				}()
				wg.Add(1)
				grpcStarted = true
			}

			grpcEndpoints <- endpoint
		}
	}

	close(grpcEndpoints)
	close(httpEndpoints)

	wg.Wait()
}
