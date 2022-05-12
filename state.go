package circuitbreaker

import (
	"time"
)

const (
	StateOpen     CircuitBreakerState = "open"
	StateClosed   CircuitBreakerState = "closed"
	StateHalfOpen CircuitBreakerState = "half-open"
)

type config struct {
	resetTimeout             time.Duration
	failureRate              float32
	numberOfCallsInHalfState int32
}

func defaultConfig() *config {
	return &config{
		resetTimeout:             60000 * time.Millisecond,
		failureRate:              0.5,
		numberOfCallsInHalfState: 5,
	}
}

type option func(*config)

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
	next(config *config) state
}

type openState struct {
	openedTime time.Time
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

func (s *openState) next(config *config) state {
	if time.Now().After(s.openedTime.Add(config.resetTimeout)) {
		return &halfOpenState{}
	}
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

func (s *closedState) next(config *config) state {
	if float32(s.failureCount)/float32(s.totalCount) > config.failureRate {
		return &openState{
			openedTime: time.Now(),
		}
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

func (s *halfOpenState) next(config *config) state {
	if s.totalCount < config.numberOfCallsInHalfState {
		return s
	} else if float32(s.failureCount)/float32(s.totalCount) > config.failureRate {
		return &openState{
			openedTime: time.Now(),
		}
	}
	return &closedState{}
}
