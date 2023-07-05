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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type requestData struct {
	Payload string
}

// Sends a HTTP POST request to the specified endpoint
func POST(service, endpoint string, port int, payload string, headers map[string]string) (int, *RestResponse, error) {
	url := fmt.Sprintf("http://%s:%d/%s", service, port, endpoint)
	postData, _ := json.Marshal(requestData{Payload: payload})
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(postData))

	headers["Content-Type"] = "application/json"
	// Forward any headers set by the user
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	endpointResponse := &RestResponse{}
	err = json.Unmarshal(data, endpointResponse)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, endpointResponse, nil
}
