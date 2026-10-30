package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/ivmerk/goshop/order/internal/api/order/v1"
	grpcClient "github.com/ivmerk/goshop/order/internal/client/grpc"
	customMiddleware "github.com/ivmerk/goshop/order/internal/middleware"
	orderRepo "github.com/ivmerk/goshop/order/internal/repository/order"
	orderSvc "github.com/ivmerk/goshop/order/internal/service/order"
	"github.com/ivmerk/goshop/platform/pkg/migrator"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
	defaultMigrationsDir = "order/migrations"
	inventoryAddr        = "localhost:50051"
	paymentAddr          = "localhost:50052"
)

func closeConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("ошибка закрытия соединения: %v", err)
	}
}

func main() {
	ctx := context.Background()

	if err := godotenv.Load("order/.env"); err != nil {
		log.Printf(".env не загружен: %v", err)
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Printf("переменная окружения DB_URI не задана (см. order/.env.example)")
		return
	}

	// Создаем соединение с базой данных
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Printf("база данных недоступна: %v", err)
		return
	}

	// Инициализируем мигратор
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = defaultMigrationsDir
	}
	db := stdlib.OpenDBFromPool(pool)
	migratorRunner := migrator.NewMigrator(db, migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(inventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ошибка подключения к inventory: %v", err)
		return
	}
	defer closeConn(inventoryConn)

	paymentConn, err := grpc.NewClient(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ошибка подключения к payment: %v", err)
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
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
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
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
