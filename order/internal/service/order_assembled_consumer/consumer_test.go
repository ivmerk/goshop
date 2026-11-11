package order_assembled_consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/order/internal/model"
	repoMocks "github.com/ivmerk/goshop/order/internal/repository/mocks"
	"github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

const (
	orderUUID = "11111111-1111-1111-1111-111111111111"
	userUUID  = "22222222-2222-2222-2222-222222222222"
	eventUUID = "33333333-3333-3333-3333-333333333333"
)

// stubConsumer запоминает переданный handler, чтобы проверить RunConsumer без Kafka.
type stubConsumer struct {
	handler consumer.MessageHandler
	err     error
}

func (c *stubConsumer) Consume(_ context.Context, h consumer.MessageHandler) error {
	c.handler = h
	return c.err
}

type ConsumerSuite struct {
	suite.Suite

	ctx             context.Context
	orderRepository *repoMocks.OrderRepository
	kafkaConsumer   *stubConsumer
	service         *service
}

func (s *ConsumerSuite) SetupSuite() {
	logger.SetNopLogger()
}

func (s *ConsumerSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.kafkaConsumer = &stubConsumer{}
	s.service = NewService(s.kafkaConsumer, s.orderRepository)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}

func (s *ConsumerSuite) message() consumer.Message {
	value, err := proto.Marshal(&eventsV1.ShipAssembled{
		EventUuid:    eventUUID,
		OrderUuid:    orderUUID,
		UserUuid:     userUUID,
		BuildTimeSec: 5,
	})
	s.Require().NoError(err)

	return consumer.Message{Key: []byte(orderUUID), Value: value, Topic: "order.assembled"}
}

func (s *ConsumerSuite) TestHandleSuccess() {
	s.orderRepository.On("Update", s.ctx, orderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusAssembled}).
		Return(orderUUID, nil)

	err := s.service.handle(s.ctx, s.message())

	s.Require().NoError(err)
}

func (s *ConsumerSuite) TestHandleOrderNotFound() {
	s.orderRepository.On("Update", s.ctx, orderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusAssembled}).
		Return("", model.ErrOrderNotFound)

	err := s.service.handle(s.ctx, s.message())

	s.Require().NoError(err) // повторять бессмысленно, сообщение пропускаем
}

func (s *ConsumerSuite) TestHandleRepositoryError() {
	updateErr := errors.New("db unavailable")
	s.orderRepository.On("Update", s.ctx, orderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusAssembled}).
		Return("", updateErr)

	err := s.service.handle(s.ctx, s.message())

	s.Require().ErrorIs(err, updateErr)
}

func (s *ConsumerSuite) TestHandleInvalidPayload() {
	// репозиторий не должен вызываться: мок без ожиданий упадёт при любом вызове
	err := s.service.handle(s.ctx, consumer.Message{Value: []byte("garbage")})

	s.Require().NoError(err) // «ядовитое» сообщение пропускаем
}

func (s *ConsumerSuite) TestRunConsumerPassesHandler() {
	s.kafkaConsumer.err = context.Canceled

	err := s.service.RunConsumer(s.ctx)

	s.Require().ErrorIs(err, context.Canceled)
	s.Require().NotNil(s.kafkaConsumer.handler)
}
