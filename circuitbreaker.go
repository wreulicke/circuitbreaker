package circuitbreaker

import (
	"errors"
	"sync"
)

type CircuitBreaker struct {
	lock   sync.RWMutex
	config *config
	state  state
}

type Hook func(old, new CircuitBreakerState)

func New(opts ...option) *CircuitBreaker {
	config := defaultConfig()

	for _, o := range opts {
		o(config)
	}

	return &CircuitBreaker{
		config: config,
		state:  &closedState{},
	}
}

func (c *CircuitBreaker) success() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.state.success()
	c.setState(c.state.next(c))
}

func (c *CircuitBreaker) failure() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.state.failure()
	c.setState(c.state.next(c))
}

func (c *CircuitBreaker) ready() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.state.state() != StateOpen
}

func (c *CircuitBreaker) setStateWithLock(to state) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.setState(to)
}

func (c *CircuitBreaker) setState(newState state) {
	old := c.state.state()
	new := newState.state()
	c.state = newState
	for _, h := range c.config.hooks {
		h(old, new)
	}
}

func (c *CircuitBreaker) do(f func() error) error {
	if !c.ready() {
		return errors.New("circuit breaker opens")
	}
	err := f()
	if err == nil {
		c.success()
	} else if c.config.ignoreError == nil || !c.config.ignoreError(err) {
		c.failure()
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
