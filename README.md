## circuit-breaker

This package provides circuit breaker
You can see concept [here](https://martinfowler.com/bliki/CircuitBreaker.html)

## Usage

```go
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
```