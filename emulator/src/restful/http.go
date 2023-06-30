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
	"strings"
	"sync"

	"encoding/json"
	"net/http"
)

const httpAddress = ":5000"

type restResponse struct {
	Status   string `json:"status"`
	Endpoint string `json:"endpoint,omitempty"`
}

// Send a response of type application/json
func writeJSONResponse(status int, response any, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	encoder := json.NewEncoder(writer)
	encoder.Encode(response)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		writeJSONResponse(http.StatusOK, &restResponse{Status: "ok"}, writer)
	} else {
		notFoundHandler(writer, request)
	}
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	endpoint := strings.TrimPrefix(request.URL.Path, "/")
	response := &restResponse{Status: "not-found", Endpoint: endpoint}

	writeJSONResponse(http.StatusNotFound, response, writer)
}

func endpointHandler(writer http.ResponseWriter, request *http.Request, endpoint *model.Endpoint) {
	response := &restResponse{Status: "ok", Endpoint: endpoint.Name}
	writeJSONResponse(http.StatusOK, response, writer)
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

	err := http.ListenAndServe(httpAddress, mux)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
