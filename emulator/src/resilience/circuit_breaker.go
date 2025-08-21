package resilience

import (
	model "application-model"
	"net/http"
	"sync"
)

type CircuitBreakerState int

const (
	OPEN CircuitBreakerState = iota
	CLOSED
	HALF_OPEN
)

type CircuitBreaker interface {
	ProxyHTTP(request *http.Request) (*http.Response, error)
	ProxyGRPC()
}

type CircuitBreakerImpl struct {
	State   CircuitBreakerState
	Timeout int // timeout in seconds

}
type CircuitBreakerRegister struct {
	ProtectedEndpoints map[string]*CircuitBreakerImpl
}

func (cbr *CircuitBreakerRegister) RegisterEndpoint(endpoint string, config *model.CircuitBreakerConfig) {
	cbr.ProtectedEndpoints[endpoint] = &CircuitBreakerImpl{
		State:   CLOSED,
		Timeout: config.Timeout,
	}
}

func (cbr *CircuitBreakerRegister) GetCircuitBreaker(endpoint string) *CircuitBreakerImpl {
	circuitBreaker, ok := cbr.ProtectedEndpoints[endpoint]
	if !ok {
		return nil
	}
	return circuitBreaker
}

var lock = &sync.Mutex{}
var circuitBreakerInstance *CircuitBreakerRegister

func CheckCircuitBreakerConfig(endpoint *model.Endpoint) bool {
	return endpoint.ResiliencePatterns.CircuitBreaker != nil
}

func GetCircuitBreakerRegister() *CircuitBreakerRegister {
	if circuitBreakerInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if circuitBreakerInstance == nil {
			circuitBreakerInstance = &CircuitBreakerRegister{
				ProtectedEndpoints: make(map[string]*CircuitBreakerImpl),
			}
		}
	}
	return circuitBreakerInstance
}
