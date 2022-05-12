package circuitbreaker

import (
	"time"

	"github.com/benbjohnson/clock"
)

const (
	StateOpen     CircuitBreakerState = "open"
	StateClosed   CircuitBreakerState = "closed"
	StateHalfOpen CircuitBreakerState = "half-open"
)

type config struct {
	clock                    clock.Clock
	resetTimeout             time.Duration
	failureRate              float32
	numberOfCallsInHalfState int32
}

func defaultConfig() *config {
	return &config{
		clock:                    clock.New(),
		resetTimeout:             60000 * time.Millisecond,
		failureRate:              0.5,
		numberOfCallsInHalfState: 5,
	}
}

type option func(*config)

func Clock(clock clock.Clock) option {
	return func(c *config) {
		c.clock = clock
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

func NumberOfCallsInHalfState(n int32) option {
	return func(c *config) {
		c.numberOfCallsInHalfState = n
	}
}

type CircuitBreakerState string

type state interface {
	state() CircuitBreakerState
	success()
	failure()
	next(config *CircuitBreaker) state
}

type openState struct {
}

type closedState struct {
	failureCount int32
	totalCount   int32
}

type halfOpenState struct {
	failureCount int32
	totalCount   int32
}

func (*openState) state() CircuitBreakerState {
	return StateOpen
}

func (s *openState) success() {
	panic("should not reache here")
}

func (s *openState) failure() {
	panic("should not reache here")
}

func (s *openState) next(c *CircuitBreaker) state {
	// do not transition to half open here, only in timer
	return s
}

func (*closedState) state() CircuitBreakerState {
	return StateClosed
}

func (s *closedState) success() {
	s.totalCount++
}

func (s *closedState) failure() {
	s.failureCount++
	s.totalCount++
}

func (s *closedState) next(c *CircuitBreaker) state {
	if float32(s.failureCount)/float32(s.totalCount) > c.config.failureRate {
		return newOpenState(c)
	}
	return s
}

func (*halfOpenState) state() CircuitBreakerState {
	return StateHalfOpen
}

func (s *halfOpenState) success() {
	s.totalCount++
}

func (s *halfOpenState) failure() {
	s.totalCount++
	s.failureCount++
}

func (s *halfOpenState) next(c *CircuitBreaker) state {
	if s.totalCount < c.config.numberOfCallsInHalfState {
		return s
	} else if float32(s.failureCount)/float32(s.totalCount) > c.config.failureRate {
		return newOpenState(c)
	}
	return &closedState{}
}

func newOpenState(c *CircuitBreaker) *openState {
	s := &openState{}
	c.config.clock.AfterFunc(c.config.resetTimeout, func() {
		c.setState(&halfOpenState{})
	})
	return s
}
