package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/wreulicke/circuitbreaker"
)

func TestExample(t *testing.T) {
	c := circuitbreaker.New(
		circuitbreaker.ResetTimeout(30*time.Second),
		circuitbreaker.FailureRate(0.5),
	)

	circuitbreaker.GuardBy(c, func() (string, error) {
		return "", nil
	})
}