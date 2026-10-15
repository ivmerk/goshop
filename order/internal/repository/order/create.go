package order

import (
	"context"
	"time"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order *model.Order) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.UUID == "" {
		return "", model.ErrInvalidOrderUUID
	}

	if _, exists := r.orders[order.UUID]; exists {
		return "", model.ErrOrderAlreadyExists
	}

	repoOrder := repoConverter.ToRepositoryOrder(order)
	repoOrder.CreatedAt = time.Now()
	r.orders[repoOrder.UUID] = repoOrder

	return repoOrder.UUID, nil
}
