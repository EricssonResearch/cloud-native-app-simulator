package circuit_breaker

import (
	"application-emulator/src/generated/client"
	"application-model/generated"
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerState string

const (
	OPEN      CircuitBreakerState = "OPEN"
	CLOSED    CircuitBreakerState = "CLOSED"
	HALF_OPEN CircuitBreakerState = "HALF OPEN"
)

var GRPC_ERROR = status.Error(codes.Unavailable, "Service unavailable")
var HTTP_ERROR = errors.New("Service unavailable")

type RequestCallback func(ctx context.Context) (any, error)

type CircuitBreaker interface {
	ProxyHTTP(request *http.Request) (*http.Response, error)
	ProxyGRPC(conn *grpc.ClientConn, service, endpoint string, request *generated.Request, options ...grpc.CallOption) (*generated.Response, error)
	ProcessRequest(cb RequestCallback, requestError error) (any, error)
}

type CircuitBreakerImpl struct {
	State             CircuitBreakerState
	Timeout           int // timeout in seconds
	RetryTimer        int // In how many seconds should we retry
	lock              sync.Mutex
	EndpointProtected string
}

func (c *CircuitBreakerImpl) ProcessRequest(cb RequestCallback, requestError error) (any, error) {
	log.Printf("[CIRCUIT BREAKER] Circuit breaker of %s in state %s\n", c.EndpointProtected, c.State)
	if c.State == OPEN {
		return nil, requestError
	}
	if c.State == HALF_OPEN {
		c.lock.Lock()
		c.State = OPEN
		c.lock.Unlock()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()

	response, err := cb(ctx)
	log.Printf("[CIRCUIT BREAKER] Request sended and callback returned %v, with error %v\n", response, err)

	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		log.Printf("[CIRCUIT BREAKER] Circuit breaker of %s timedout\n", c.EndpointProtected)
		c.lock.Lock()
		c.State = OPEN
		c.lock.Unlock()
		go func() {
			time.Sleep(time.Second * time.Duration(c.RetryTimer))
			c.lock.Lock()
			c.State = HALF_OPEN
			c.lock.Unlock()
			return
		}()
		return nil, requestError
	}

	if c.State == OPEN || c.State == HALF_OPEN {
		c.lock.Lock()
		c.State = CLOSED
		c.lock.Unlock()
	}
	return response, err
}

func (c *CircuitBreakerImpl) ProxyHTTP(request *http.Request) (*http.Response, error) {

	response, err := c.ProcessRequest(func(ctx context.Context) (any, error) {
		request = request.WithContext(ctx)
		response, err := http.DefaultClient.Do(request)
		return response, err
	}, HTTP_ERROR)

	if err != nil {
		return nil, err
	}

	httpResponse, ok := response.(*http.Response)
	log.Printf("[CIRCUIT BREAKER] The response returned from circuit breaker as %v", httpResponse)

	if !ok {
		return nil, errors.New("HTTP response from Circuit breaker callback broken")
	}

	return httpResponse, err
}

func (c *CircuitBreakerImpl) ProxyGRPC(conn *grpc.ClientConn, service, endpoint string, request *generated.Request, options ...grpc.CallOption) (*generated.Response, error) {
	response, err := c.ProcessRequest(func(ctx context.Context) (any, error) {
		return client.CallGeneratedEndpoint(ctx, conn, service, endpoint, request, options...)
	}, GRPC_ERROR)

	if err != nil {
		return nil, err
	}
	grpcResponse, ok := response.(*generated.Response)
	if !ok {
		return nil, errors.New("GRPC response from Circuit breaker callback broken")
	}

	return grpcResponse, err
}
