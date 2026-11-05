package order_assembled_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/assembly/internal/model"
	"github.com/ivmerk/goshop/platform/pkg/kafka"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

type service struct {
	producer kafka.Producer
}

func NewService(producer kafka.Producer) *service {
	return &service{producer: producer}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	return s.producer.Send(ctx, []byte(event.OrderUUID), payload)
}
