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

package restful

import (
	"cloud-native-app-simulator/model"
	"fmt"
	"net/http"
	"sync"
)

const httpPort = 5000

// Launches a HTTP server to serve one or more endpoints
func HTTP(endpointChannel chan model.Endpoint, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	endpoints := []model.Endpoint{}
	for endpoint := range endpointChannel {
		endpoints = append(endpoints, endpoint)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprint(writer, "{\"status\": \"ok\"}\n")
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
