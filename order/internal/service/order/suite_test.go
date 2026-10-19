package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/ivmerk/goshop/order/internal/client/grpc/mocks"
	repoMocks "github.com/ivmerk/goshop/order/internal/repository/mocks"
)

const (
	orderUUID = "11111111-1111-1111-1111-111111111111"
	userUUID  = "22222222-2222-2222-2222-222222222222"
	partUUID1 = "33333333-3333-3333-3333-333333333333"
	partUUID2 = "44444444-4444-4444-4444-444444444444"
	txUUID    = "55555555-5555-5555-5555-555555555555"
)

type ServiceSuite struct {
	suite.Suite

	ctx             context.Context
	orderRepository *repoMocks.OrderRepository
	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient
	service         *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())
	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
