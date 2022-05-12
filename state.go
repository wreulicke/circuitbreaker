package circuitbreaker

const (
	StateOpen     CircuitBreakerState = "open"
	StateClosed   CircuitBreakerState = "closed"
	StateHalfOpen CircuitBreakerState = "half-open"
)

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
	// nop
}

func (s *openState) failure() {
	// nop
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
		c.setStateWithLock(&halfOpenState{})
	})
	return s
}
