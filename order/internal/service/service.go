package service

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, order *model.Order) (string, error)
	Get(ctx context.Context, uuid string) (*model.Order, error)
	Pay(ctx context.Context, uuid string, transactionID string) (string, error)
	Cancel(ctx context.Context, uuid string) (string, error)
}
