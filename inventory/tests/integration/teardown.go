//go:build integration

package integration

import (
	"context"
	"log"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			log.Printf("failed to terminate app: %v", err)
		}
	}
	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			log.Printf("failed to terminate mongo: %v", err)
		}
	}
	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			log.Printf("failed to remove network: %v", err)
		}
	}
}
