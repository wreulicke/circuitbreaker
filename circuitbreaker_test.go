package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker(t *testing.T) {
	c := New()

	GuardBy(c, func() (string, error) {
		return "", nil
	})
	GuardBy(c, func() (string, error) {
		return "", errors.New("test")
	})
	GuardBy(c, func() (string, error) {
		return "", errors.New("test")
	})
	assert.Equal(t, StateOpen, c.state.state())
}

func TestCircuitBreakerOpen(t *testing.T) {
	c := New(ResetTimeout(1 * time.Second))
	c.state = &openState{
		openedTime: time.Now().Add(-3 * time.Second),
	}

	GuardBy(c, func() (string, error) {
		return "", nil
	})
	assert.Equal(t, StateHalfOpen, c.state.state())
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	c := New()
	c.state = &halfOpenState{}

	GuardBy(c, func() (string, error) {
		return "", nil
	})
	assert.Equal(t, StateClosed, c.state.state())
}

func TestCircuitBreakerHalfOpen_Failure(t *testing.T) {
	c := New()
	c.state = &halfOpenState{}

	GuardBy(c, func() (string, error) {
		return "", errors.New("test")
	})
	assert.Equal(t, StateOpen, c.state.state())
}
