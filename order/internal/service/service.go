package service

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	Get(ctx context.Context, uuid string) (*model.Order, error)
	Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, uuid string) error
}
