package middleware

import (
	"net/http"
	"realtime-bidding-system/pkg/resilience"

	"github.com/gin-gonic/gin"
)

// Bulkhead returns a middleware that limits concurrent requests for a specific route or group
func Bulkhead(b resilience.Bulkhead) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := b.Execute(c.Request.Context(), func() error {
			c.Next()
			return nil
		})

		if err != nil {
			if err == resilience.ErrBulkheadFull {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "bulkhead limit reached",
				})
				return
			}
			// If context is cancelled
			c.AbortWithStatus(http.StatusRequestTimeout)
			return
		}
	}
}
