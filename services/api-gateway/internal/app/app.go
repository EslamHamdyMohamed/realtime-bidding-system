package app

import (
	"context"
	"net/http"
	"realtime-bidding-system/pkg/config"
	"realtime-bidding-system/pkg/logger"
	"realtime-bidding-system/pkg/postgres"
	"realtime-bidding-system/pkg/redis"
	"realtime-bidding-system/pkg/tracing"
	gh "realtime-bidding-system/services/api-gateway/internal/http"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	pg     *postgres.DB
	rdb    *redis.Client
	tp     *sdktrace.TracerProvider
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	// Initialize Tracing
	tp, err := tracing.Init(ctx, cfg.ServiceName, cfg.OTLPURL)
	if err != nil {
		logger.Log.Warn("Failed to initialize tracing", zap.Error(err))
	}

	// Initialize Dependencies
	var pg *postgres.DB
	if cfg.PostgresURL != "" {
		pg, err = postgres.New(cfg.PostgresURL)
		if err != nil {
			logger.Log.Warn("Failed to connect to Postgres", zap.Error(err))
		}
	}

	var rdb *redis.Client
	if cfg.RedisURL != "" {
		rdb = redis.New(cfg.RedisURL)
	}

	handlers := gh.NewHandlers(pg, rdb)
	router := gh.NewRouter(handlers, cfg.ServiceName)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &App{
		cfg:    cfg,
		server: server,
		pg:     pg,
		rdb:    rdb,
		tp:     tp,
	}, nil
}

func (a *App) Run() error {
	logger.Log.Info("Starting API Gateway (Gin)", zap.String("port", a.cfg.HTTPPort))
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	logger.Log.Info("Shutting down API Gateway")
	if a.tp != nil {
		if err := a.tp.Shutdown(ctx); err != nil {
			logger.Log.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}
	return a.server.Shutdown(ctx)
}
