package service

import (
	"context"

	"github.com/ivmerk/goshop/assembly/internal/model"
)

// ConsumerService читает OrderPaid и запускает сборку.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// ProducerService публикует ShipAssembled.
type ProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
