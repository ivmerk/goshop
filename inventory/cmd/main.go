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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryAPI "github.com/ivmerk/goshop/inventory/internal/api/inventory/v1"
	"github.com/ivmerk/goshop/inventory/internal/config"
	partRepo "github.com/ivmerk/goshop/inventory/internal/repository/part"
	partService "github.com/ivmerk/goshop/inventory/internal/service/part"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

const configPath = "deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %w", err)
		return
	}

	dbURI := config.AppConfig().Mongo.URI()
	mongoDBName := config.AppConfig().Mongo.DatabaseName()

	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("ошибка подключения к mongo: %v", err)
		return
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("ошибка отключения от mongo: %v", err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("база данных недоступна: %v", err)
		return
	}
	collection := client.Database(mongoDBName).Collection("parts")

	repo := partRepo.NewPartRepository(collection)
	if err := repo.Seed(ctx); err != nil {
		log.Printf("не удалось заполнить коллекцию: %v", err)
		return
	}

	svc := partService.NewPartService(repo)
	api := inventoryAPI.NewAPI(svc)

	server := grpc.NewServer()
	inventoryV1.RegisterInventoryServiceServer(server, api)

	reflection.Register(server)

	grpcAddr := config.AppConfig().GRPC.Address()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	go func() {
		log.Printf("🚀 gRPC InventoryService запущен на порту %s", grpcAddr)
		if err := server.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы gRPC-сервера...")
	server.GracefulStop()
	log.Println("✅ Сервер остановлен")
}
