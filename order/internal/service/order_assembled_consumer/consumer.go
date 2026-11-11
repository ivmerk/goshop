package order_assembled_consumer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ivmerk/goshop/order/internal/converter/kafka/decoder"
	"github.com/ivmerk/goshop/order/internal/model"
	"github.com/ivmerk/goshop/order/internal/repository"
	"github.com/ivmerk/goshop/platform/pkg/kafka"
	"github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

type service struct {
	consumer        kafka.Consumer
	orderRepository repository.OrderRepository
}

func NewService(c kafka.Consumer, repo repository.OrderRepository) *service {
	return &service{consumer: c, orderRepository: repo}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderAssembled consumer")
	return s.consumer.Consume(ctx, s.handle)
}

func (s *service) handle(ctx context.Context, msg consumer.Message) error {
	event, err := decoder.ShipAssembledToModel(msg.Value)
	if err != nil {
		logger.Error(ctx, "failed to decode ShipAssembled", zap.Error(err))
		return nil // «ядовитое» сообщение пропускаем, ретрай ничего не изменит
	}
	logger.Info(ctx, "ShipAssembled received", zap.String("order_uuid", event.OrderUUID))

	_, err = s.orderRepository.Update(ctx, event.OrderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusAssembled})
	if errors.Is(err, model.ErrOrderNotFound) {
		logger.Warn(ctx, "order not found for ShipAssembled", zap.String("order_uuid", event.OrderUUID))
		return nil
	}

	return err
}
