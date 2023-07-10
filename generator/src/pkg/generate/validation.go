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
	s "application-generator/src/pkg/service"
	model "application-model"

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

			if endpoint.NetworkComplexity != nil {
				for _, calledService := range endpoint.NetworkComplexity.CalledServices {
					errs = validation.IsDNS1035Label(calledService.Service)

					if len(errs) > 0 {
						return fmt.Errorf("Call from endpoint '%s' to invalid service '%s': %s", endpoint.Name, calledService.Service, errs[0])
					}

					errs = validation.IsDNS1123Subdomain(calledService.Endpoint)

					if len(errs) > 0 {
						return fmt.Errorf("Call from endpoint '%s' to invalid endpoint '%s': %s", endpoint.Name, calledService.Endpoint, errs[0])
					}
				}
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
			service.Resources.Limits.Cpu,
			service.Resources.Limits.Memory,
			service.Resources.Requests.Cpu,
			service.Resources.Requests.Memory,
		}

		for _, limit := range limits {
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

	return nil
}

// Validate that protocols are set in both endpoint definition and call
func ValidateProtocols(service *model.Service, endpoint *model.Endpoint) error {
	validProtocols := map[string]bool{"http": true, "grpc": true}

	if !validProtocols[endpoint.Protocol] {
		return fmt.Errorf("Endpoint '%s' in service '%s' has invalid protocol '%s'",
			endpoint.Name, service.Name, endpoint.Protocol)
	}

	if endpoint.NetworkComplexity != nil {
		for _, calledService := range endpoint.NetworkComplexity.CalledServices {
			if !validProtocols[calledService.Protocol] {
				return fmt.Errorf("Call to endpoint '%s' from endpoint '%s' has invalid protocol '%s'",
					calledService.Endpoint, endpoint.Name, calledService.Protocol)
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
			return fmt.Errorf("Service '%s' needs to be deployed on at least one cluster", service.Name)
		}

		if len(service.Endpoints) == 0 {
			return fmt.Errorf("At least one endpoint is required in service '%s'", service.Name)
		} else {
			for _, endpoint := range service.Endpoints {
				err := ValidateProtocols(&service, &endpoint)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Validates an input JSON config provided by the user
func ValidateFileConfig(config *model.FileConfig) error {
	if err := ValidateRequiredParameters(config); err != nil {
		return err
	}

	if err := ValidateNames(config); err != nil {
		return err
	}

	if err := ValidateResources(config); err != nil {
		return err
	}

	return nil
}

// Applies default values to input JSON
func ApplyDefaults(config *model.FileConfig) {
	for i := range config.Services {
		service := &config.Services[i]

		if service.Resources.Limits.Cpu == "" {
			service.Resources.Limits.Cpu = s.LimitsCPUDefault
		}
		if service.Resources.Requests.Cpu == "" {
			service.Resources.Requests.Cpu = s.RequestsCPUDefault
		}
		if service.Resources.Limits.Memory == "" {
			service.Resources.Limits.Memory = s.LimitsMemoryDefault
		}
		if service.Resources.Requests.Memory == "" {
			service.Resources.Requests.Memory = s.RequestsMemoryDefault
		}

		if service.Processes <= 0 {
			service.Processes = s.SvcProcessesDefault
		}

		if service.ReadinessProbe <= 0 {
			service.ReadinessProbe = s.SvcReadinessProbeDefault
		}

		for j := range service.Clusters {
			cluster := &service.Clusters[j]

			if cluster.Namespace == "" {
				cluster.Namespace = s.ClusterNamespaceDefault
			}
		}

		for k := range service.Endpoints {
			endpoint := &service.Endpoints[k]

			if endpoint.CpuComplexity != nil && endpoint.CpuComplexity.Threads < 1 {
				endpoint.CpuComplexity.Threads = 1
			}
		}
	}
}
