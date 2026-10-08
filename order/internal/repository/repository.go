package repository

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	DeleteOrder(ctx context.Context, uuid string) error
}
