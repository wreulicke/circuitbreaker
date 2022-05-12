package circuitbreaker

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrOpen = errors.New("circuit breaker is open")
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
		return ErrOpen
	}
	err := f()
	if err == nil || c.config.isSuccessful(err) {
		c.success()
	} else if !c.config.isIgnorable(err) {
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

type IgnorableError struct {
	err error
}

func (e *IgnorableError) Error() string {
	return fmt.Sprintf("marked as a ignorable error: %s", e.err.Error())
}

func (e *IgnorableError) Unwrap() error {
	return e.err
}

func Ignore(err error) error {
	return &IgnorableError{err: err}
}

type SuccessfulError struct {
	err error
}

func (e *SuccessfulError) Error() string {
	return fmt.Sprintf("marked as a successful error: %s", e.err.Error())
}

func (e *SuccessfulError) Unwrap() error {
	return e.err
}

func Success(err error) error {
	return &SuccessfulError{
		err: err,
	}
}
