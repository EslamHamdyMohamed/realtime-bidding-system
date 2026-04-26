package middleware

import (
	"net/http"
	"realtime-bidding-system/pkg/resilience"

	"github.com/gin-gonic/gin"
)

// CircuitBreaker returns a middleware that wraps the handler in a circuit breaker
func CircuitBreaker(cb resilience.CircuitBreaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := cb.Execute(func() (interface{}, error) {
			c.Next()
			// If any other middleware or handler aborted with error, we should report it to CB
			if len(c.Errors) > 0 {
				return nil, c.Errors.Last()
			}
			return nil, nil
		})

		if err != nil {
			// Check if the error is from the circuit breaker being open
			// gobreaker doesn't export the error types directly in a way that matches easily without casting
			// but we can check the error message or assume any error from Execute that isn't from c.Next is a CB error.

			// If the handler didn't abort but CB return error, it means CB is open
			if !c.IsAborted() {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"error": "circuit breaker is open",
				})
			}
		}
	}
}
