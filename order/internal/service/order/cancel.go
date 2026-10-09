package order

import (
	"context"
	"time"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
)

func (s *service) Cancel(ctx context.Context, uuid string) (string, error) {

	order := &model.Order{UUID: uuid}
	repoOrder := repoConverter.ToRepositoryOrder(order)
	repoOrder.Transaction = time.Now().Format("20060102150405") // Пример генерации идентификатора транзакции на основе текущего времени

	return repoOrder.Transaction, nil
}
