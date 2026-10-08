package order

import (
	"context"
	"time"

	"github.com/ivmerk/goshop/order/internal/model"
	repoModel "github.com/ivmerk/goshop/order/internal/repository/model"
	"github.com/samber/lo"
)

func (r *repository) Update(_ context.Context, uuid string, updatedOrder *model.OrderUpdateInfo) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingOrder, ok := r.orders[uuid]
	if !ok {
		return "", model.ErrOrderNotFound
	}

	if updatedOrder.Transaction != nil {
		existingOrder.Transaction = updatedOrder.Transaction
	}
	existingOrder.Status = repoModel.OrderStatus(updatedOrder.Status)
	existingOrder.UpdatedAt = lo.ToPtr(time.Now())

	return existingOrder.UUID, nil
}
