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
	resetTimeout time.Duration
	failureRate  float32
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
	failureCount int
	totalCount   int
	failureRate  float32
}

type halfOpenState struct {
	failed bool
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
	if float32(s.failureCount)/float32(s.totalCount) > s.failureRate {
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
}

func (s *halfOpenState) failure() {
	s.failed = true
}

func (s *halfOpenState) next(config *config) state {
	if s.failed {
		return &openState{
			openedTime: time.Now(),
		}
	}
	return &closedState{
		failureRate: config.failureRate,
	}
}
