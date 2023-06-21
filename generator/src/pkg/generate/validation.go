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

package generate

import (
	"application-generator/src/pkg/model"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation"
)

// Occurrences returns the number of occurrences of every value in a slice of strings
func Occurrences(strSlice []string) map[string]int {
	occurrences := make(map[string]int)
	for _, entry := range strSlice {
		occurrences[entry]++
	}
	return occurrences
}

// Validates service and endpoint names in JSON config
func ValidateNames(config *model.FileConfig) error {
	serviceNames := []string{}

	// Validate service names (RFC 1035 DNS Label)
	for _, service := range config.Services {
		errs := validation.IsDNS1035Label(service.Name)

		// There can be several conformance errors but only one is returned by this function
		// If the user fixes one error, the next error will be shown when running the generator again
		if len(errs) > 0 {
			return fmt.Errorf("Service '%s' has invalid name: %s", service.Name, errs[0])
		}

		serviceNames = append(serviceNames, service.Name)
	}

	serviceNameOccurrences := Occurrences(serviceNames)

	for _, service := range config.Services {
		// Duplicate name found
		if serviceNameOccurrences[service.Name] > 1 {
			return fmt.Errorf("Duplicate service name '%s'", service.Name)
		}

		endpointNames := []string{}

		// Validate endpoint names (RFC 1123 DNS Subdomain)
		for _, endpoint := range service.Endpoints {
			errs := validation.IsDNS1123Subdomain(endpoint.Name)

			if len(errs) > 0 {
				return fmt.Errorf("Endpoint '%s' has invalid name: %s", endpoint.Name, errs[0])
			}

			endpointNames = append(endpointNames, endpoint.Name)
		}

		endpointNameOccurrences := Occurrences(serviceNames)

		for _, endpoint := range service.Endpoints {
			// Duplicate name found
			if endpointNameOccurrences[endpoint.Name] > 1 {
				return fmt.Errorf("Duplicate endpoint '%s' in service '%s'", endpoint.Name, service.Name)
			}
		}
	}

	return nil
}

// Validates resource limits in input JSON
func ValidateResources(config *model.FileConfig) error {
	for _, service := range config.Services {
		limits := []string{
			service.Resources.Requests.Cpu,
			service.Resources.Requests.Memory,
			service.Resources.Limits.Cpu,
			service.Resources.Limits.Memory,
		}

		for _, limit := range limits {
			// If the user hasn't provided a request or limit, they will be set to their default values later
			if limit != "" {
				quantity, err := resource.ParseQuantity(limit)

				if err != nil {
					return fmt.Errorf("Invalid resource allocation '%s': %s", limit, err)
				}

				// TODO: Max limits
				if quantity.Sign() != 1 {
					return fmt.Errorf("Resource allocation '%s' too low", limit)
				}
			}
		}
	}

	return nil
}

// Validate that input JSON contains required parameters
func ValidateRequiredParameters(config *model.FileConfig) error {
	if len(config.Services) == 0 {
		return errors.New("At least one service is required")
	}

	for _, service := range config.Services {
		if len(service.Clusters) == 0 {
			return fmt.Errorf("Service '%s' needs to be placed on at least one cluster", service.Name)
		}

		if len(service.Endpoints) == 0 {
			return fmt.Errorf("At least one endpoint is required in service '%s'", service.Name)
		}
	}

	// TODO: Check clusters, called_services, etc

	return nil
}

// Validates an input JSON config provided by the user
func ValidateFileConfig(config *model.FileConfig) error {
	if err := ValidateNames(config); err != nil {
		return err
	}

	if err := ValidateResources(config); err != nil {
		return err
	}

	if err := ValidateRequiredParameters(config); err != nil {
		return err
	}

	return nil
}
