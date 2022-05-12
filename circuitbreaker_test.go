package circuitbreaker

import (
	"errors"
	"testing"
	"time"

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
	c := New(ResetTimeout(100 * time.Second))
	c.state = &openState{
		openedTime: time.Now(),
	}

	assertDontCalled(t, c, success)
	assert.Equal(t, StateOpen, c.state.state())
}

func TestCircuitBreakerOpen_Timeout(t *testing.T) {
	c := New(ResetTimeout(1 * time.Second))
	c.state = &openState{
		openedTime: time.Now().Add(-3 * time.Second),
	}

	assertCalled(t, c, success)
	assert.Equal(t, StateHalfOpen, c.state.state())
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	c := New(NumberOfCallsInHalfState(1))
	c.state = &halfOpenState{}

	assertCalled(t, c, success)
	assert.Equal(t, StateClosed, c.state.state())
}

func TestCircuitBreakerHalfOpen_Failure(t *testing.T) {
	c := New(NumberOfCallsInHalfState(1))
	c.state = &halfOpenState{}

	assertCalled(t, c, failure)
	assert.Equal(t, StateOpen, c.state.state())
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
