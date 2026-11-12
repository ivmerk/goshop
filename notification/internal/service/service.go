package service

import (
	"context"

	"github.com/ivmerk/goshop/notification/internal/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error
}
