package order

import (
	grpcClient "github.com/ivmerk/goshop/order/internal/client/grpc"
	"github.com/ivmerk/goshop/order/internal/repository"
	def "github.com/ivmerk/goshop/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpcClient.InventoryClient,
	paymentClient grpcClient.PaymentClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
