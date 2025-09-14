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

package client

import (
	"application-emulator/src/resilience/circuit_breaker"
	"application-model/generated"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

const useProtoJSON = true

// Sends a HTTP POST request to the specified endpoint
func POST(service, endpoint string, port int, payload string, headers http.Header, sourceEndpoint string) (int, *generated.Response, error) {
	var url string
	// Omit the port if zero
	if port == 0 {
		url = fmt.Sprintf("http://%s/%s", service, endpoint)
	} else {
		url = fmt.Sprintf("http://%s:%d/%s", service, port, endpoint)
	}

	var postData []byte

	if useProtoJSON {
		marshalOptions := protojson.MarshalOptions{UseProtoNames: true, AllowPartial: true}
		postData, _ = marshalOptions.Marshal(&generated.Request{Payload: payload})
	} else {
		postData, _ = json.Marshal(&generated.Request{Payload: payload})
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(postData))

	// Forward any other headers set by the user
	for key, values := range headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	// Send the request
	// Here it comes the circuit breaker

	circuitBreakerRegistry := circuit_breaker.GetCircuitBreakerRegistry()
	cbName := circuitBreakerRegistry.BuildName(sourceEndpoint, service, endpoint)
	circuitBreaker := circuitBreakerRegistry.GetCircuitBreaker(cbName)

	log.Printf("[CLIENT HTTP] Circuit breaker obtanied %v", circuitBreaker)

	var response *http.Response
	var err error

	if circuitBreaker == nil {
		response, err = http.DefaultClient.Do(request)
	} else {
		response, err = circuitBreaker.ProxyHTTP(request)
	}

	log.Printf("[CLIENT HTTP] Request returned from %s/%s", service, endpoint)
	if err != nil {
		return 0, nil, err
	}

	// Read all bytes
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	// Assume that we received a valid model.RESTResponse
	endpointResponse := &generated.Response{}
	if useProtoJSON {
		err = protojson.Unmarshal(data, endpointResponse)
		if err != nil {
			return 0, nil, err
		}
	} else {
		err = json.Unmarshal(data, endpointResponse)
		if err != nil {
			return 0, nil, err
		}
	}

	return response.StatusCode, endpointResponse, nil
}
