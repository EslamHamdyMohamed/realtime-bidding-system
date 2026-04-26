package middleware

import (
	"net/http"
	"realtime-bidding-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("RequestID")
				logger.Log.Error("panic recovered",
					zap.String("request_id", requestID),
					zap.Any("error", err),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
