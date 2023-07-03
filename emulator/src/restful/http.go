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
	"cloud-native-app-simulator/emulator/src/stressors"
	"cloud-native-app-simulator/emulator/src/util"
	"cloud-native-app-simulator/model"

	"fmt"
	"strings"
	"sync"

	"encoding/json"
	"net/http"
)

const httpAddress = ":5000"

type RestResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"message,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`
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
		writeJSONResponse(http.StatusOK, &RestResponse{Status: "ok"}, writer)
	} else {
		notFoundHandler(writer, request)
	}
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	endpoint := strings.TrimPrefix(request.URL.Path, "/")
	response := &RestResponse{
		Status:       "error",
		ErrorMessage: fmt.Sprintf("Endpoint %s doesn't exist", endpoint),
		Endpoint:     endpoint,
	}

	writeJSONResponse(http.StatusNotFound, response, writer)
}

type endpointHandler struct {
	endpoint *model.Endpoint
}

func (handler *endpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	trace := util.TraceEndpointStart(handler.endpoint)
	response := &RestResponse{Status: "ok", Endpoint: handler.endpoint.Name}

	// TODO: Async (goroutines)
	if handler.endpoint.CpuComplexity != nil {
		stressors.CPU(handler.endpoint.CpuComplexity)
	}

	writeJSONResponse(http.StatusOK, response, writer)
	util.TraceEndpointEnd(trace)
}

// Launch a HTTP server to serve one or more endpoints
func HTTP(endpointChannel chan model.Endpoint, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	for endpoint := range endpointChannel {
		mux.Handle(fmt.Sprintf("/%s", endpoint.Name), &endpointHandler{endpoint: &endpoint})
	}

	err := http.ListenAndServe(httpAddress, mux)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
