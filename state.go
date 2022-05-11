package circuitbreaker

import "time"

const (
	Open = iota
	Closed
	HalfOpen
)

type config struct {
	resetTimeout time.Duration
	failureRate  float32
}

type CircuitBreakerState int

type state interface {
	state() CircuitBreakerState
	success()
	failure()
	next(config *config) state
}

type openState struct {
	openedTime   time.Time
	resetTimeout time.Duration
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
	return Open
}

func (s *openState) success() {
	panic("should not reache here")
}

func (s *openState) failure() {
	panic("should not reache here")
}

func (s *openState) next(config *config) state {
	if s.openedTime.Add(s.resetTimeout).After(time.Now()) {
		return &halfOpenState{}
	}
	return s
}

func (*closedState) state() CircuitBreakerState {
	return Closed
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
			openedTime:   time.Now(),
			resetTimeout: config.resetTimeout,
		}
	}
	return s
}

func (*halfOpenState) state() CircuitBreakerState {
	return HalfOpen
}

func (s *halfOpenState) success() {
}

func (s *halfOpenState) failure() {
	s.failed = true
}

func (s *halfOpenState) next(config *config) state {
	if s.failed {
		return &openState{
			openedTime:   time.Now(),
			resetTimeout: config.resetTimeout,
		}
	}
	return &closedState{
		failureRate: config.failureRate,
	}
}
