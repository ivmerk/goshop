package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (*model.Order, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return &model.Order{}, model.ErrOrderNotFound
	}
	return repoConverter.FromRepositoryOrder(order), nil

}
