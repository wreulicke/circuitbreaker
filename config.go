package circuitbreaker

import (
	"time"

	"github.com/benbjohnson/clock"
)

type config struct {
	clock clock.Clock

	resetTimeout             time.Duration
	failureRate              float32
	numberOfCallsInHalfState int32
	hooks                    []Hook
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

func WithClock(clock clock.Clock) option {
	return func(c *config) {
		c.clock = clock
	}
}

func WithResetTimeout(d time.Duration) option {
	return func(c *config) {
		c.resetTimeout = d
	}
}

func WithFailureRate(p float32) option {
	return func(c *config) {
		c.failureRate = p
	}
}

func WithNumberOfCallsInHalfState(n int32) option {
	return func(c *config) {
		c.numberOfCallsInHalfState = n
	}
}

func WithHook(h Hook) option {
	return func(c *config) {
		c.hooks = append(c.hooks, h)
	}
}
