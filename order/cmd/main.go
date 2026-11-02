package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/ivmerk/goshop/order/internal/api/order/v1"
	grpcClient "github.com/ivmerk/goshop/order/internal/client/grpc"
	"github.com/ivmerk/goshop/order/internal/config"
	customMiddleware "github.com/ivmerk/goshop/order/internal/middleware"
	orderRepo "github.com/ivmerk/goshop/order/internal/repository/order"
	orderSvc "github.com/ivmerk/goshop/order/internal/service/order"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	"github.com/ivmerk/goshop/platform/pkg/migrator"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

const shutdownTimeout = 10 * time.Second

func closeConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("ошибка закрытия соединения: %v", err)
	}
}

func main() {
	ctx := context.Background()

	const configPath = "deploy/compose/order/.env"
	if err := config.Load(configPath); err != nil {
		log.Printf("ошибка загрузки конфигурации: %v", err)
		return
	}

	cfg := config.AppConfig()

	if err := logger.Init(cfg.Logger.Level(), cfg.Logger.AsJson()); err != nil {
		log.Printf("failed to init logger: %v", err)
		return
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	// Создаем соединение с базой данных
	pool, err := pgxpool.New(ctx, cfg.Postgres.URI())
	if err != nil {
		logger.Error(ctx, "failed to connect to database", zap.Error(err))
		return
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error(ctx, "failed to ping database", zap.Error(err))
		return
	}

	// Инициализируем мигратор
	db := stdlib.OpenDBFromPool(pool)
	migratorRunner := migrator.NewMigrator(db, cfg.Postgres.MigrationsDir())

	err = migratorRunner.Up()
	if err != nil {
		logger.Error(ctx, "failed to migrate database", zap.Error(err))
		return
	}

	inventoryConn, err := grpc.NewClient(cfg.InventoryGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(ctx, "failed to connect to inventory", zap.Error(err))
		return
	}
	defer closeConn(inventoryConn)

	paymentConn, err := grpc.NewClient(cfg.PaymentGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(ctx, "failed to connect to payment", zap.Error(err))
		return
	}
	defer closeConn(paymentConn)

	orderRepository := orderRepo.NewOrderRepository(pool)

	inventoryClient := grpcClient.NewInventoryClient(inventoryV1.NewInventoryServiceClient(inventoryConn))
	paymentClient := grpcClient.NewPaymentClient(paymentV1.NewPaymentServiceClient(paymentConn))

	orderService := orderSvc.NewService(orderRepository, inventoryClient, paymentClient)

	orderHandler := orderAPI.NewOrderHandler(orderService)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		logger.Error(ctx, "failed to create OpenAPI server", zap.Error(err))
		return
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              cfg.OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: cfg.OrderHTTP.ReadTimeout(),
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		logger.Info(ctx, "starting HTTP server", zap.String("address", cfg.OrderHTTP.Address()))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "failed to start HTTP server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "shutting down HTTP server")

	// Создаем контекст с таймаутом для остановки сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error(ctx, "failed to shutdown HTTP server", zap.Error(err))
	}

	logger.Info(ctx, "graceful shutdown complete")
}
