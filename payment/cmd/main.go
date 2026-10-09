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
	paymentSvc "github.com/ivmerk/goshop/payment/internal/service/payment"

	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

const grpcPort = "50052"

func main() {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := paymentSvc.NewService()
	api := paymentAPI.NewPaymentAPI(svc)
	server := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(server, api)

	// Позволяет grpcurl видеть список сервисов без .proto-файла.
	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC PaymentService запущен на порту %s", grpcPort)
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
