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
	"application-emulator/src/resilience"
	"application-emulator/src/stressors"
	"application-emulator/src/util"
	model "application-model"
	"application-model/generated"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
)

const useProtoJSON = true

// Send a response of type application/json
func writeJSONResponse(status int, response *generated.Response, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if useProtoJSON {
		marshalOptions := protojson.MarshalOptions{UseProtoNames: true, AllowPartial: true}
		data, err := marshalOptions.Marshal(response)
		if err != nil {
			panic(err)
		}

		writer.Write(data)
		writer.Write([]byte{'\n'})
	} else {
		data, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}

		writer.Write(data)
		writer.Write([]byte{'\n'})
	}
}

// This should always be available because the readiness probe will send a request to this address
func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		writeJSONResponse(http.StatusOK, &generated.Response{}, writer)
	} else {
		notFoundHandler(writer, request)
	}
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	endpoint := strings.TrimPrefix(request.URL.Path, "/")
	response := &generated.Response{
		Message:  fmt.Sprintf("Endpoint %s doesn't exist", endpoint),
		Endpoint: endpoint,
	}

	writeJSONResponse(http.StatusNotFound, response, writer)
}

type endpointHandler struct {
	endpoint *model.Endpoint
}

func (handler endpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	trace := util.TraceEndpointCall(handler.endpoint, "HTTP")
	response := &generated.Response{
		Endpoint: handler.endpoint.Name,
		Tasks:    stressors.Exec(request, handler.endpoint),
	}
	writeJSONResponse(http.StatusOK, response, writer)
	util.LogEndpointCall(trace)
}

// Launch a HTTP server to serve one or more endpoints
func HTTP(endpoints []model.Endpoint) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	circuitBreakerInstance := resilience.GetCircuitBreakerRegister()
	for i := range endpoints {
		if resilience.CheckCircuitBreakerConfig(&endpoints[i]) {
			circuitBreakerInstance.RegisterEndpoint(endpoints[i].Name, endpoints[i].ResiliencePatterns.CircuitBreaker)
		}
		mux.Handle(fmt.Sprintf("/%s", endpoints[i].Name), endpointHandler{endpoint: &endpoints[i]})
	}

	err := http.ListenAndServe(":5000", mux)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
