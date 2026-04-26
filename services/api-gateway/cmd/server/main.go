package main

import (
	"context"
	"os"
	"os/signal"
	"realtime-bidding-system/pkg/config"
	"realtime-bidding-system/pkg/logger"
	"realtime-bidding-system/services/api-gateway/internal/app"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.ServiceName)
	defer logger.Sync()

	application, err := app.New(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to initialize app", zap.Error(err))
	}

	go func() {
		if err := application.Run(); err != nil {
			logger.Log.Fatal("Failed to run app", zap.Error(err))
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	logger.Log.Info("Shutdown signal received", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		logger.Log.Error("App shutdown failed", zap.Error(err))
	} else {
		logger.Log.Info("App shutdown completed")
	}
}
