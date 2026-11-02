package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentAPI "github.com/ivmerk/goshop/payment/internal/api/payment/v1"
	"github.com/ivmerk/goshop/payment/internal/config"
	paymentSvc "github.com/ivmerk/goshop/payment/internal/service/payment"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

const configPath = "deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return
	}

	grpcAddr := config.AppConfig().GRPC.Address()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	svc := paymentSvc.NewService()
	api := paymentAPI.NewPaymentAPI(svc)
	server := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(server, api)

	// Позволяет grpcurl видеть список сервисов без .proto-файла.
	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC PaymentService запущен на порту %s", grpcAddr)
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
