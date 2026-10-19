package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *service) Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error) {
	if method == model.PaymentMethodUnknown {
		return "", model.ErrInvalidPaymentMethod
	}
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return "", err
	}
	if order.Status != model.OrderStatusPendingPayment {
		return "", model.ErrOrderNotPayable
	}

	txID, err := s.paymentClient.PayOrder(ctx, order.UUID, order.User, method)
	if err != nil {
		return "", err
	}

	_, err = s.orderRepository.Update(ctx, uuid, &model.OrderUpdateInfo{
		Transaction: &txID, Status: model.OrderStatusPaid, Payment: method,
	})
	return txID, err
}
