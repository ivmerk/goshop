package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
)

func (s *service) Pay(ctx context.Context, uuid string, transactionID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := &model.Order{UUID: uuid}
	repoOrder := repoConverter.ToRepositoryOrder(order)
	repoOrder.Transaction = transactionID

	return repoOrder.Transaction, nil
}
