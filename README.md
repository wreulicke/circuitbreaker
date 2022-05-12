## circuit-breaker

This package provides circuit breaker.
You can see the concept [here](https://martinfowler.com/bliki/CircuitBreaker.html).

## Usage

```go
import (
	"testing"
	"time"

	"github.com/wreulicke/circuitbreaker"
)

func TestExample(t *testing.T) {
	// initialize with default options
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

```