package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/ivmerk/goshop/notification/internal/converter/kafka/decoder"
	def "github.com/ivmerk/goshop/notification/internal/service"
	"github.com/ivmerk/goshop/platform/pkg/kafka"
	"github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

type service struct {
	consumer        kafka.Consumer
	telegramService def.TelegramService
}

func NewService(c kafka.Consumer, telegramService def.TelegramService) *service {
	return &service{consumer: c, telegramService: telegramService}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaid consumer")

	return s.consumer.Consume(ctx, s.handle)
}

func (s *service) handle(ctx context.Context, msg consumer.Message) error {
	event, err := decoder.OrderPaidToModel(msg.Value)
	if err != nil {
		logger.Error(ctx, "failed to decode OrderPaid", zap.Error(err))
		return nil // битое сообщение повторно не обработать: пропускаем
	}

	logger.Info(ctx, "OrderPaid received", zap.String("order_uuid", event.OrderUUID))

	return s.telegramService.SendOrderPaidNotification(ctx, event)
}
