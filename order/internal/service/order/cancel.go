package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {

	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return err
	}
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return nil
	}

	updatedOrder := &model.OrderUpdateInfo{
		Status: model.OrderStatusCancelled,
	}
	_, err = s.orderRepository.Update(ctx, order.UUID, updatedOrder)
	return err
}
