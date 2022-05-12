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

func (c *CircuitBreaker) ready() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.state.state() != StateOpen
}

func (c *CircuitBreaker) do(f func() error) error {
	if !c.ready() {
		c.next()
		return errors.New("circuit breaker opens")
	}
	err := f()
	if err != nil {
		c.failure()
	} else {
		c.success()
	}
	return err
}

func GuardBy[T any](cb *CircuitBreaker, f func() (T, error)) (r T, err error) {
	err = cb.do(func() error {
		r, err = f()
		return err
	})
	return r, err
}
