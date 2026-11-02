package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentAPI "github.com/ivmerk/goshop/payment/internal/api/payment/v1"
	"github.com/ivmerk/goshop/payment/internal/config"
	paymentSvc "github.com/ivmerk/goshop/payment/internal/service/payment"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

const configPath = "deploy/compose/payment/.env"

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

	grpcAddr := cfg.GRPC.Address()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Error(ctx, "failed to listen", zap.Error(err))
		return
	}

	svc := paymentSvc.NewService()
	api := paymentAPI.NewPaymentAPI(svc)
	server := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(server, api)

	// Позволяет grpcurl видеть список сервисов без .proto-файла.
	reflection.Register(server)

	go func() {
		logger.Info(ctx, "🚀 gRPC PaymentService started", zap.String("address", grpcAddr))
		if err := server.Serve(lis); err != nil {
			logger.Error(ctx, "failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "👋 gRPC PaymentService stopping")
	server.GracefulStop()
	logger.Info(ctx, "👋 gRPC PaymentService stopped")
}
