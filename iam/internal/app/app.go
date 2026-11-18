package app

import (
	"context"
	"errors"
	"net"

	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ivmerk/goshop/iam/internal/config"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/grpc/health"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	"github.com/ivmerk/goshop/platform/pkg/migrator"
	authV1 "github.com/ivmerk/goshop/shared/pkg/proto/auth/v1"
	userV1 "github.com/ivmerk/goshop/shared/pkg/proto/user/v1"
)

type App struct {
	di         *diContainer
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	logger.Info(ctx, "🚀 gRPC IAM server started", zap.String("address", a.listener.Addr().String()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	for _, f := range []func(context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.initMigrations,
		a.initListener,
		a.initGRPCServer,
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

// initMigrations накатывает SQL-миграции пользователей при старте сервиса.
func (a *App) initMigrations(ctx context.Context) error {
	pool, err := a.di.PostgresPool(ctx)
	if err != nil {
		return err
	}
	db := stdlib.OpenDBFromPool(pool)
	return migrator.NewMigrator(db, config.AppConfig().Postgres.MigrationsDir()).Up()
}

func (a *App) initListener(_ context.Context) error {
	lis, err := net.Listen("tcp", config.AppConfig().IAMGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(context.Context) error {
		if err := lis.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
	a.listener = lis
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	authAPI, err := a.di.AuthV1API(ctx)
	if err != nil {
		return err
	}
	userAPI, err := a.di.UserV1API(ctx)
	if err != nil {
		return err
	}

	a.grpcServer = grpc.NewServer()
	closer.AddNamed("gRPC server", func(context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)

	authV1.RegisterAuthServiceServer(a.grpcServer, authAPI)
	userV1.RegisterUserServiceServer(a.grpcServer, userAPI)

	return nil
}
