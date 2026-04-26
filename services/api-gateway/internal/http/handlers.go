package http

import (
	"net/http"
	"realtime-bidding-system/pkg/postgres"
	"realtime-bidding-system/pkg/redis"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	pg  *postgres.DB
	rdb *redis.Client
}

func NewHandlers(pg *postgres.DB, rdb *redis.Client) *Handlers {
	return &Handlers{pg: pg, rdb: rdb}
}

func (h *Handlers) Health(c *gin.Context) {
	status := "OK"
	details := make(map[string]string)

	// Check Postgres
	if h.pg != nil {
		if err := h.pg.Pool.Ping(c.Request.Context()); err != nil {
			status = "Degraded"
			details["postgres"] = "Down"
		} else {
			details["postgres"] = "Up"
		}
	}

	// Check Redis
	if h.rdb != nil {
		if err := h.rdb.Ping(c.Request.Context()); err != nil {
			status = "Degraded"
			details["redis"] = "Down"
		} else {
			details["redis"] = "Up"
		}
	}

	res := gin.H{
		"status":  status,
		"details": details,
	}

	if status == "OK" {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusServiceUnavailable, res)
	}
}
