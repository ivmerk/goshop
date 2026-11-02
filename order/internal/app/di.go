package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/ivmerk/goshop/order/internal/api/order/v1"
	grpcClient "github.com/ivmerk/goshop/order/internal/client/grpc"
	"github.com/ivmerk/goshop/order/internal/config"
	"github.com/ivmerk/goshop/order/internal/repository"
	orderRepo "github.com/ivmerk/goshop/order/internal/repository/order"
	"github.com/ivmerk/goshop/order/internal/service"
	orderSvc "github.com/ivmerk/goshop/order/internal/service/order"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderServer *orderV1.Server
	orderAPI    *orderAPI.OrderHandler
	orderSvc    service.OrderService
	orderRepo   repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient

	pgPool *pgxpool.Pool
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) OrderServer(ctx context.Context) (*orderV1.Server, error) {
	if d.orderServer == nil {
		h, err := d.OrderAPI(ctx)
		if err != nil {
			return nil, err
		}
		srv, err := orderV1.NewServer(h)
		if err != nil {
			return nil, err
		}
		d.orderServer = srv
	}
	return d.orderServer, nil
}

func (d *diContainer) OrderAPI(ctx context.Context) (*orderAPI.OrderHandler, error) {
	if d.orderAPI == nil {
		svc, err := d.OrderService(ctx)
		if err != nil {
			return nil, err
		}
		d.orderAPI = orderAPI.NewOrderHandler(svc)
	}
	return d.orderAPI, nil
}

func (d *diContainer) OrderService(ctx context.Context) (service.OrderService, error) {
	if d.orderSvc == nil {
		repo, err := d.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}
		inv, err := d.InventoryClient()
		if err != nil {
			return nil, err
		}
		pay, err := d.PaymentClient()
		if err != nil {
			return nil, err
		}
		d.orderSvc = orderSvc.NewService(repo, inv, pay)
	}
	return d.orderSvc, nil
}

func (d *diContainer) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if d.orderRepo == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}
		d.orderRepo = orderRepo.NewOrderRepository(pool)
	}
	return d.orderRepo, nil
}

func (d *diContainer) InventoryClient() (grpcClient.InventoryClient, error) {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		closer.AddNamed("Inventory gRPC connection", func(context.Context) error { return conn.Close() })
		d.inventoryClient = grpcClient.NewInventoryClient(inventoryV1.NewInventoryServiceClient(conn))
	}
	return d.inventoryClient, nil
}

func (d *diContainer) PaymentClient() (grpcClient.PaymentClient, error) {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		closer.AddNamed("Payment gRPC connection", func(context.Context) error { return conn.Close() })
		d.paymentClient = grpcClient.NewPaymentClient(paymentV1.NewPaymentServiceClient(conn))
	}
	return d.paymentClient, nil
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
