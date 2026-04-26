package http

import (
	"realtime-bidding-system/pkg/resilience"
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

	// Initialize Resilience Patterns
	// Bulkhead limiting global concurrent requests to 1000
	globalBulkhead := resilience.NewBulkhead(1000)
	// Circuit Breaker for the API Gateway
	apiCB := resilience.NewCircuitBreaker("api-gateway")

	// Apply Middlewares
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Tracing(serviceName))
	r.Use(middleware.Logging())
	r.Use(middleware.RateLimit(limiter))
	r.Use(middleware.Metrics())

	// Apply Resilience Middlewares
	r.Use(middleware.Bulkhead(globalBulkhead))
	r.Use(middleware.CircuitBreaker(apiCB))

	r.GET("/health", h.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
