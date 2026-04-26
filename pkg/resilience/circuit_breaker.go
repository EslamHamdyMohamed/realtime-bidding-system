package resilience

import (
	"time"

	"github.com/sony/gobreaker/v2"
)

// CircuitBreaker defines the interface for circuit breaker operations
type CircuitBreaker interface {
	Execute(func() (interface{}, error)) (interface{}, error)
}

type circuitBreaker struct {
	cb *gobreaker.CircuitBreaker[interface{}]
}

// NewCircuitBreaker creates a new circuit breaker with default settings
func NewCircuitBreaker(name string) CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	return &circuitBreaker{
		cb: gobreaker.NewCircuitBreaker[interface{}](settings),
	}
}

func (c *circuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	return c.cb.Execute(req)
}
