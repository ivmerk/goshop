package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryAPI "github.com/ivmerk/goshop/inventory/internal/api/inventory/v1"
	partRepo "github.com/ivmerk/goshop/inventory/internal/repository/part"
	partService "github.com/ivmerk/goshop/inventory/internal/service/part"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = "50051"

func main() {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	repo := partRepo.NewPartRepository()
	svc := partService.NewPartService(repo)
	api := inventoryAPI.NewAPI(svc)

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
