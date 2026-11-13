package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ivmerk/goshop/order/internal/api/health"
	"github.com/ivmerk/goshop/order/internal/config"
	customMiddleware "github.com/ivmerk/goshop/order/internal/middleware"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	"github.com/ivmerk/goshop/platform/pkg/migrator"
)

type App struct {
	di         *diContainer
	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	svc, err := a.di.OrderAssembledConsumerService(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	closer.AddNamed("consumer loop", func(context.Context) error {
		cancel()
		return nil
	})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info(ctx, "🚀 HTTP server started", zap.String("address", a.httpServer.Addr))
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := svc.RunConsumer(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, stop := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer stop()
		if err := a.httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return g.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	for _, f := range []func(context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.initMigrations,
		a.initHTTPServer,
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

func (a *App) initMigrations(ctx context.Context) error {
	pool, err := a.di.PostgresPool(ctx)
	if err != nil {
		return err
	}
	db := stdlib.OpenDBFromPool(pool)
	return migrator.NewMigrator(db, config.AppConfig().Postgres.MigrationsDir()).Up()
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := a.di.OrderServer(ctx)
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Get("/health", health.Handler)
	r.Mount("/", orderServer)

	cfg := config.AppConfig().OrderHTTP
	a.httpServer = &http.Server{
		Addr:              cfg.Address(),
		Handler:           r,
		ReadHeaderTimeout: cfg.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
	return nil
}
