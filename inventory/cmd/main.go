package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryAPI "github.com/ivmerk/goshop/inventory/internal/api/inventory/v1"
	"github.com/ivmerk/goshop/inventory/internal/config"
	partRepo "github.com/ivmerk/goshop/inventory/internal/repository/part"
	partService "github.com/ivmerk/goshop/inventory/internal/service/part"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

const configPath = "deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v", err)
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

	dbURI := cfg.Mongo.URI()
	mongoDBName := cfg.Mongo.DatabaseName()

	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err != nil {
		logger.Error(ctx, "ошибка подключения к mongo:", zap.Error(err))
		return
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.Error(ctx, "ошибка отключения от mongo:", zap.Error(err))
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error(ctx, "база данных недоступна:", zap.Error(err))
		return
	}
	collection := client.Database(mongoDBName).Collection("parts")

	repo := partRepo.NewPartRepository(collection)
	if err := repo.Seed(ctx); err != nil {
		logger.Error(ctx, "не удалось заполнить коллекцию:", zap.Error(err))
		return
	}

	svc := partService.NewPartService(repo)
	api := inventoryAPI.NewAPI(svc)

	server := grpc.NewServer()
	inventoryV1.RegisterInventoryServiceServer(server, api)

	reflection.Register(server)

	grpcAddr := cfg.GRPC.Address()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Error(ctx, "ошибка создания tcp-сокета:", zap.Error(err))
		return
	}
	go func() {
		logger.Info(ctx, "🚀 gRPC InventoryService запущен", zap.String("grpcAddr", grpcAddr))
		if err := server.Serve(lis); err != nil {
			logger.Error(ctx, "ошибка запуска gRPC-сервера:", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "🛑 Завершение работы gRPC-сервера...")
	server.GracefulStop()
	logger.Info(ctx, "✅ Сервер остановлен")
}
