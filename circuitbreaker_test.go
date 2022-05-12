package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func success() (string, error) {
	return "", nil
}

func failure() (string, error) {
	return "", errors.New("test")
}

func TestCircuitBreaker(t *testing.T) {
	c := New()

	assertCalled(t, c, success)
	assertCalled(t, c, failure)
	assertCalled(t, c, failure)
	assert.Equal(t, StateOpen, c.state.state())
}

func TestCircuitBreakerOpen(t *testing.T) {
	c := New(WithResetTimeout(100 * time.Second))
	c.state = newOpenState(c)

	assertDontCalled(t, c, success)
	assert.Equal(t, StateOpen, c.state.state())
}

func TestCircuitBreakerOpen_Timeout(t *testing.T) {
	mock := clock.NewMock()
	c := New(WithResetTimeout(1*time.Second), WithClock(mock))
	c.state = newOpenState(c)
	mock.Add(3 * time.Second)

	assertCalled(t, c, success)
	assert.Equal(t, StateHalfOpen, c.state.state())
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	c := New(WithNumberOfCallsInHalfState(1))
	c.state = &halfOpenState{}

	assertCalled(t, c, success)
	assert.Equal(t, StateClosed, c.state.state())
}

func TestCircuitBreakerHalfOpen_Failure(t *testing.T) {
	c := New(WithNumberOfCallsInHalfState(1))
	c.state = &halfOpenState{}

	assertCalled(t, c, failure)
	assert.Equal(t, StateOpen, c.state.state())
}

func TestCircuitBreakerHook(t *testing.T) {
	var old, new CircuitBreakerState
	c := New(WithNumberOfCallsInHalfState(1), WithHook(func(o, n CircuitBreakerState) {
		old, new = o, n
	}))
	c.state = &halfOpenState{}

	assertCalled(t, c, failure)
	assert.Equal(t, StateOpen, c.state.state())
	assert.Equal(t, old, StateHalfOpen)
	assert.Equal(t, new, StateOpen)
}

func assertCalled[T any](t *testing.T, cb *CircuitBreaker, f func() (T, error)) {
	t.Helper()

	var ok bool
	GuardBy(cb, func() (T, error) {
		ok = true
		return f()
	})
	assert.True(t, ok, "should be called")
}

func assertDontCalled[T any](t *testing.T, cb *CircuitBreaker, f func() (T, error)) {
	t.Helper()

	ok := true
	GuardBy(cb, func() (T, error) {
		ok = false
		return f()
	})
	assert.True(t, ok, "should not be called")
}
