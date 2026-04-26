package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(serviceName string) {
	logger, _ := zap.NewProduction()
	Log = logger.With(
		zap.String("service", serviceName),
	)
}

func Sync() {
	_ = Log.Sync()
}
