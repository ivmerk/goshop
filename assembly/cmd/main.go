package main

import (
	"context"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ivmerk/goshop/assembly/internal/app"
	"github.com/ivmerk/goshop/assembly/internal/config"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

const (
	configPath      = "deploy/compose/assembly/.env"
	shutdownTimeout = 10 * time.Second
)

func main() {
	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		logger.Error(ctx, "failed to load config", zap.Error(err))
		return
	}

	a, err := app.New(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create app", zap.Error(err))
		closeAll()
		return
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	if err := a.Run(ctx); err != nil {
		logger.Error(ctx, "app stopped with error", zap.Error(err))
	}

	closeAll()
}

func closeAll() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "failed to close resources", zap.Error(err))
	}
}
