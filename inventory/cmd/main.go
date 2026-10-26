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
	partRepo "github.com/ivmerk/goshop/inventory/internal/repository/part"
	partService "github.com/ivmerk/goshop/inventory/internal/service/part"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
	"github.com/joho/godotenv"
)

const grpcPort = "50051"

func main() {
	ctx := context.Background()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	if err := godotenv.Load("inventory/.env"); err != nil {
		log.Printf(".env не загружен: %v", err)
	}

	dbURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if dbURI == "" {
		log.Printf("переменная окружения MONGO_URI не задана (см. inventory/.env.example)")
		return
	}

	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("ошибка подключения к mongo: %v", err)
		return
	}
	defer client.Disconnect(ctx)
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

	go func() {
		log.Printf("🚀 gRPC InventoryService запущен на порту %s", grpcPort)
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
