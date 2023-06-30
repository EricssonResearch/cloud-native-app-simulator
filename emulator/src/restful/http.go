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
	"strings"

	"fmt"
	"sync"

	"net/http"
)

const httpPort = 5000

type restResponse struct {
	Status   string `json:"status"`
	Endpoint string `json:"endpoint,omitempty"`
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

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)

	endpoint := strings.TrimPrefix(request.URL.Path, "/")

	encoder := json.NewEncoder(writer)
	response := &restResponse{Status: "not found", Endpoint: endpoint}

	encoder.Encode(response)
}

func endpointHandler(writer http.ResponseWriter, request *http.Request, endpoint *model.Endpoint) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(writer)
	response := &restResponse{Status: "ok", Endpoint: endpoint.Name}

	encoder.Encode(response)
}

// Launch a HTTP server to serve one or more endpoints
func HTTP(endpointChannel chan model.Endpoint, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	for endpoint := range endpointChannel {
		mux.HandleFunc(fmt.Sprintf("/%s", endpoint.Name), func(w http.ResponseWriter, r *http.Request) {
			endpointHandler(w, r, &endpoint)
		})
	}

	address := fmt.Sprintf(":%d", httpPort)
	err := http.ListenAndServe(address, mux)

	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
