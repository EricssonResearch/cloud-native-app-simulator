package circuit_breaker

import (
	model "application-model"
	"fmt"
	"sync"
)

var instanceLock = &sync.Mutex{}
var registerLock = &sync.Mutex{}
var GetLock = &sync.Mutex{}

type CircuitBreakerRegister struct {
	ProtectedEndpoints map[string]*CircuitBreakerImpl
}

var circuitBreakerInstance *CircuitBreakerRegister

func (cbr *CircuitBreakerRegister) BuildName(sourceEndpoint, destService, destEndpoint string) string {
	return fmt.Sprintf("%s:%s/%s", sourceEndpoint, destService, destEndpoint)
}
func (cbr *CircuitBreakerRegister) RegisterEndpoint(endpoint string, config *model.CircuitBreakerConfig) {
	registerLock.Lock()
	defer registerLock.Unlock()
	cbr.ProtectedEndpoints[endpoint] = &CircuitBreakerImpl{
		State:             CLOSED,
		Timeout:           config.Timeout,
		RetryTimer:        config.RetryTimer,
		EndpointProtected: endpoint,
	}
}

func (cbr *CircuitBreakerRegister) GetCircuitBreaker(endpoint string) *CircuitBreakerImpl {
	GetLock.Lock()
	defer GetLock.Unlock()

	circuitBreaker, ok := cbr.ProtectedEndpoints[endpoint]
	if !ok {
		return nil
	}
	return circuitBreaker
}

func CheckCircuitBreakerConfig(endpoint *model.Endpoint) bool {
	return endpoint.ResiliencePatterns.CircuitBreaker != nil
}

func GetCircuitBreakerRegistry() *CircuitBreakerRegister {
	if circuitBreakerInstance == nil {
		instanceLock.Lock()
		defer instanceLock.Unlock()
		if circuitBreakerInstance == nil {
			circuitBreakerInstance = &CircuitBreakerRegister{
				ProtectedEndpoints: make(map[string]*CircuitBreakerImpl),
			}
		}
	}
	return circuitBreakerInstance
}
