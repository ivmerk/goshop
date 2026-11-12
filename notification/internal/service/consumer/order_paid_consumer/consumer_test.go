package order_paid_consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/notification/internal/model"
	serviceMocks "github.com/ivmerk/goshop/notification/internal/service/mocks"
	"github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

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
	telegramService *serviceMocks.TelegramService
	kafkaConsumer   *stubConsumer
	service         *service
}

func (s *ConsumerSuite) SetupSuite() {
	logger.SetNopLogger()
}

func (s *ConsumerSuite) SetupTest() {
	s.ctx = context.Background()
	s.telegramService = serviceMocks.NewTelegramService(s.T())
	s.kafkaConsumer = &stubConsumer{}
	s.service = NewService(s.kafkaConsumer, s.telegramService)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}

func (s *ConsumerSuite) message() consumer.Message {
	value, err := proto.Marshal(&eventsV1.OrderPaid{
		EventUuid:       "event-1",
		OrderUuid:       "order-1",
		UserUuid:        "user-1",
		PaymentMethod:   "CARD",
		TransactionUuid: "tx-1",
	})
	s.Require().NoError(err)

	return consumer.Message{Value: value, Topic: "order.paid"}
}

func (s *ConsumerSuite) TestHandleSuccess() {
	s.telegramService.On("SendOrderPaidNotification", s.ctx, model.OrderPaidEvent{
		EventUUID:       "event-1",
		OrderUUID:       "order-1",
		UserUUID:        "user-1",
		PaymentMethod:   "CARD",
		TransactionUUID: "tx-1",
	}).Return(nil)

	s.Require().NoError(s.service.handle(s.ctx, s.message()))
}

func (s *ConsumerSuite) TestHandleTelegramError() {
	sendErr := errors.New("telegram unavailable")
	s.telegramService.On("SendOrderPaidNotification", s.ctx, model.OrderPaidEvent{
		EventUUID:       "event-1",
		OrderUUID:       "order-1",
		UserUUID:        "user-1",
		PaymentMethod:   "CARD",
		TransactionUUID: "tx-1",
	}).Return(sendErr)

	s.Require().ErrorIs(s.service.handle(s.ctx, s.message()), sendErr)
}

func (s *ConsumerSuite) TestHandleInvalidPayload() {
	err := s.service.handle(s.ctx, consumer.Message{Value: []byte("garbage")})

	s.Require().NoError(err) // сервис Telegram не вызывается, сообщение пропускаем
}

func (s *ConsumerSuite) TestRunConsumerPassesHandler() {
	s.kafkaConsumer.err = context.Canceled

	err := s.service.RunConsumer(s.ctx)

	s.Require().ErrorIs(err, context.Canceled)
	s.Require().NotNil(s.kafkaConsumer.handler)
}
