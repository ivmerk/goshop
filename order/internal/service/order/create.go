package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *service) Create(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		return nil, model.ErrEmptyParts
	}
	prices, err := s.inventoryClient.ListPartPrices(ctx, partUUIDs)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, id := range partUUIDs {
		price, ok := prices[id]
		if !ok {
			return nil, fmt.Errorf("%w: %s", model.ErrPartNotFound, id)
		}
		total += price
	}

	order := &model.Order{
		UUID: uuid.NewString(), User: userUUID, Parts: partUUIDs,
		Total: total, Status: model.OrderStatusPendingPayment,
	}
	if _, err := s.orderRepository.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil

}
