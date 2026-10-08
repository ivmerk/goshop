package repository

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (string, error)
	Get(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, uuid string, updatedOrder *model.OrderUpdateInfo) (string, error)
	Delete(ctx context.Context, uuid string) error
}
