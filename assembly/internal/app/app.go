package app

import (
	"context"
	"errors"

	"github.com/ivmerk/goshop/assembly/internal/config"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

type App struct {
	di *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	svc, err := a.di.ConsumerService()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	closer.AddNamed("consumer loop", func(context.Context) error {
		cancel()
		return nil
	})

	if err := svc.RunConsumer(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	for _, f := range []func(context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initDI,
	} {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := config.AppConfig().Logger
	return logger.Init(cfg.Level(), cfg.AsJson())
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	closer.AddNamed("logger", func(context.Context) error { return logger.Sync() })
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.di = NewDiContainer()
	return nil
}
