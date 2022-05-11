package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type CircuitBreaker struct {
	lock   sync.RWMutex
	config *config
	state  state
}

func defaulConfig() *config {
	return &config{
		resetTimeout: 60000 * time.Millisecond,
		failureRate:  0.5,
	}
}

type option func(*config)

func New(opts ...option) *CircuitBreaker {
	config := defaulConfig()

	for _, o := range opts {
		o(config)
	}

	return &CircuitBreaker{
		config: config,
		state: &closedState{
			failureRate: config.failureRate,
		},
	}
}

func ResetTimeout(d time.Duration) option {
	return func(c *config) {
		c.resetTimeout = d
	}
}

func FailureRate(p float32) option {
	return func(c *config) {
		c.failureRate = p
	}
}

func (c *CircuitBreaker) success() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.state.success()
	c.state = c.state.next(c.config)
}

func (c *CircuitBreaker) failure() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.state.failure()
	c.state = c.state.next(c.config)
}

func (c *CircuitBreaker) next() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.state = c.state.next(c.config)
}

func GuardBy[T any](cb *CircuitBreaker, f func() (T, error)) (r T, err error) {
	switch cb.state.state() {
	case StateOpen:
		cb.next()
		return r, errors.New("circuit breaker opens")
	case StateHalfOpen:
		r, err = f()
		if err != nil {
			cb.failure()
		} else {
			cb.success()
		}
		return r, err
	case StateClosed:
		r, err = f()
		if err != nil {
			cb.failure()
		} else {
			cb.success()
		}
		return r, err
	default:
		panic("should reach here")
	}
}
