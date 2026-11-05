package order_paid_consumer

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ivmerk/goshop/assembly/internal/converter/kafka/decoder"
	"github.com/ivmerk/goshop/assembly/internal/model"
	def "github.com/ivmerk/goshop/assembly/internal/service"
	"github.com/ivmerk/goshop/platform/pkg/kafka"
	"github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

type service struct {
	consumer kafka.Consumer
	producer def.ProducerService
}

func NewService(c kafka.Consumer, p def.ProducerService) *service {
	return &service{consumer: c, producer: p}
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

	n, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		return err
	}
	delay := time.Duration(1+n.Int64()) * time.Second
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	return s.producer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(delay / time.Second),
	})
}
