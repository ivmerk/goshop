package app

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	authAPI "github.com/ivmerk/goshop/iam/internal/api/auth/v1"
	userAPI "github.com/ivmerk/goshop/iam/internal/api/user/v1"
	"github.com/ivmerk/goshop/iam/internal/config"
	"github.com/ivmerk/goshop/iam/internal/repository"
	sessionRepo "github.com/ivmerk/goshop/iam/internal/repository/session"
	userRepo "github.com/ivmerk/goshop/iam/internal/repository/user"
	"github.com/ivmerk/goshop/iam/internal/service"
	authSvc "github.com/ivmerk/goshop/iam/internal/service/auth"
	userSvc "github.com/ivmerk/goshop/iam/internal/service/user"
	"github.com/ivmerk/goshop/platform/pkg/cache"
	redisClient "github.com/ivmerk/goshop/platform/pkg/cache/redis"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	authV1 "github.com/ivmerk/goshop/shared/pkg/proto/auth/v1"
	userV1 "github.com/ivmerk/goshop/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API authV1.AuthServiceServer
	userV1API userV1.UserServiceServer

	authService service.AuthService
	userService service.UserService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	pgPool      *pgxpool.Pool
	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) AuthV1API(ctx context.Context) (authV1.AuthServiceServer, error) {
	if d.authV1API == nil {
		svc, err := d.AuthService(ctx)
		if err != nil {
			return nil, err
		}
		d.authV1API = authAPI.NewAPI(svc)
	}
	return d.authV1API, nil
}

func (d *diContainer) UserV1API(ctx context.Context) (userV1.UserServiceServer, error) {
	if d.userV1API == nil {
		svc, err := d.UserService(ctx)
		if err != nil {
			return nil, err
		}
		d.userV1API = userAPI.NewAPI(svc)
	}
	return d.userV1API, nil
}

func (d *diContainer) AuthService(ctx context.Context) (service.AuthService, error) {
	if d.authService == nil {
		users, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		sessions, err := d.SessionRepository(ctx)
		if err != nil {
			return nil, err
		}
		d.authService = authSvc.NewService(users, sessions, config.AppConfig().Session.TTL())
	}
	return d.authService, nil
}

func (d *diContainer) UserService(ctx context.Context) (service.UserService, error) {
	if d.userService == nil {
		users, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		d.userService = userSvc.NewService(users)
	}
	return d.userService, nil
}

func (d *diContainer) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if d.userRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}
		d.userRepository = userRepo.NewUserRepository(pool)
	}
	return d.userRepository, nil
}

func (d *diContainer) SessionRepository(ctx context.Context) (repository.SessionRepository, error) {
	if d.sessionRepository == nil {
		client, err := d.RedisClient(ctx)
		if err != nil {
			return nil, err
		}
		d.sessionRepository = sessionRepo.NewSessionRepository(client)
	}
	return d.sessionRepository, nil
}

func (d *diContainer) PostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			return nil, err
		}
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, err
		}
		closer.AddNamed("PostgreSQL pool", func(context.Context) error {
			pool.Close()
			return nil
		})
		d.pgPool = pool
	}
	return d.pgPool, nil
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		cfg := config.AppConfig().Redis
		pool := &redigo.Pool{
			MaxIdle:     cfg.MaxIdle(),
			IdleTimeout: cfg.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", cfg.Address())
			},
		}
		closer.AddNamed("Redis pool", func(context.Context) error {
			return pool.Close()
		})
		d.redisPool = pool
	}
	return d.redisPool
}

func (d *diContainer) RedisClient(ctx context.Context) (cache.RedisClient, error) {
	if d.redisClient == nil {
		client := redisClient.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
		if err := client.Ping(ctx); err != nil {
			return nil, err
		}
		d.redisClient = client
	}
	return d.redisClient, nil
}
