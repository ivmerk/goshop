package order

import (
	"context"
	"time"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order *model.Order) (string, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := repoConverter.ToRepositoryOrder(order)
	repoOrder.CreatedAt = time.Now()

	return repoOrder.UUID, nil

}
