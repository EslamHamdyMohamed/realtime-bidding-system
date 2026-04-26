package http

import (
	"realtime-bidding-system/services/api-gateway/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

func NewRouter(h *Handlers, serviceName string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Initialize Rate Limiter
	limiter := middleware.NewIPRateLimiter(rate.Limit(100), 200)

	// Apply Middlewares
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Tracing(serviceName))
	r.Use(middleware.Logging())
	r.Use(middleware.RateLimit(limiter))
	r.Use(middleware.Metrics())

	r.GET("/health", h.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
