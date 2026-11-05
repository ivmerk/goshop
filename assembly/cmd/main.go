package main

import (
	"context"
	"log"
	"syscall"
	"time"

	"github.com/ivmerk/goshop/assembly/internal/app"
	"github.com/ivmerk/goshop/assembly/internal/config"
	"github.com/ivmerk/goshop/platform/pkg/closer"
)

const (
	configPath      = "deploy/compose/assembly/.env"
	shutdownTimeout = 10 * time.Second
)

func main() {
	if err := config.Load(configPath); err != nil {
		log.Printf("failed to load config: %v", err)
		return
	}

	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Printf("failed to create app: %v", err)
		closeAll()
		return
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	if err := a.Run(ctx); err != nil {
		log.Printf("app stopped with error: %v", err)
	}

	closeAll()
}

func closeAll() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		log.Printf("failed to close resources: %v", err)
	}
}
