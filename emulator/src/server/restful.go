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

package server

import (
	"application-emulator/src/stressors"
	"application-emulator/src/util"
	model "application-model"

	"fmt"
	"strings"
	"sync"

	"encoding/json"
	"net/http"
)

const httpAddress = ":5000"

// Send a response of type application/json
func writeJSONResponse(status int, response any, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	encoder := json.NewEncoder(writer)
	encoder.Encode(response)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		writeJSONResponse(http.StatusOK, &model.RESTResponse{Status: "ok"}, writer)
	} else {
		notFoundHandler(writer, request)
	}
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	endpoint := strings.TrimPrefix(request.URL.Path, "/")
	response := &model.RESTResponse{
		Status:       "error",
		ErrorMessage: fmt.Sprintf("Endpoint %s doesn't exist", endpoint),
		Endpoint:     endpoint,
	}

	writeJSONResponse(http.StatusNotFound, response, writer)
}

type endpointHandler struct {
	endpoint *model.Endpoint
}

func (handler endpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	trace := util.TraceEndpointCall(handler.endpoint)
	response := &model.RESTResponse{Status: "ok", Endpoint: handler.endpoint.Name}

	if handler.endpoint.ExecutionMode == "parallel" {
		response.TaskResponses = stressors.ExecParallel(request, handler.endpoint)
	} else {
		response.TaskResponses = stressors.ExecSequential(request, handler.endpoint)
	}

	writeJSONResponse(http.StatusOK, response, writer)
	util.PrintEndpointTrace(trace)
}

// Launch a HTTP server to serve one or more endpoints
func HTTP(endpointChannel chan model.Endpoint, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	for endpoint := range endpointChannel {
		mux.Handle(fmt.Sprintf("/%s", endpoint.Name), endpointHandler{endpoint: &endpoint})
	}

	err := http.ListenAndServe(httpAddress, mux)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
