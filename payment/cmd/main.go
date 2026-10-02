package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

const grpcPort = "50052"

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUUID}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(server, &paymentService{})

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
