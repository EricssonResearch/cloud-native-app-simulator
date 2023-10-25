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
	"application-model/generated"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

const useProtoJSON = true

// Sends a HTTP POST request to the specified endpoint
func POST(service, endpoint string, port int, payload string, headers http.Header) (int, *generated.Response, error) {
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
	response, err := http.DefaultClient.Do(request)
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
