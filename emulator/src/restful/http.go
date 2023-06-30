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
	"encoding/json"

	"fmt"
	"sync"

	"net/http"
)

const httpPort = 5000

type restResponse struct {
	Status   string `json:"status"`
	Endpoint string `json:"endpoint,omitempty"`
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)

	encoder := json.NewEncoder(writer)
	response := &restResponse{Status: "not found", Endpoint: request.URL.Path}

	encoder.Encode(response)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(writer)
		response := &restResponse{Status: "ok"}

		encoder.Encode(response)
	} else {
		notFoundHandler(writer, request)
	}
}

// Launches a HTTP server to serve one or more endpoints
func HTTP(endpointChannel chan model.Endpoint, wg *sync.WaitGroup) {
	defer wg.Done()

	endpoints := []model.Endpoint{}
	for endpoint := range endpointChannel {
		endpoints = append(endpoints, endpoint)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	address := fmt.Sprintf(":%d", httpPort)
	err := http.ListenAndServe(address, mux)

	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
