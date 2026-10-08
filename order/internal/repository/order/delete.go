package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (r *repository) Delete(_ context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	delete(r.orders, uuid)
	return nil
}
