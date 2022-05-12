package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/wreulicke/circuitbreaker"
)

func TestExample(t *testing.T) {
	c := circuitbreaker.New(
		circuitbreaker.WithResetTimeout(60*time.Second),
		circuitbreaker.WithFailureRate(0.5),
		circuitbreaker.WithNumberOfCallsInHalfState(5),
		circuitbreaker.WithIsIgnorable(circuitbreaker.DefaultIsIgnorable),
		circuitbreaker.WithIsSuccessful(circuitbreaker.DefaultIsSuccessful),
	)

	circuitbreaker.GuardBy(c, func() (string, error) {
		return "", nil
	})
}
