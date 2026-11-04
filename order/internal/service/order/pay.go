package order

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ivmerk/goshop/order/internal/model"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

func (s *service) Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	if method == model.PaymentMethodUnknown {
		return "", model.ErrInvalidPaymentMethod
	}
	order, err := s.orderRepository.Get(ctx, orderUUID)
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

	if _, err = s.orderRepository.Update(ctx, orderUUID, &model.OrderUpdateInfo{
		Transaction: &txID, Status: model.OrderStatusPaid, Payment: method,
	}); err != nil {
		return "", err
	}

	if err := s.orderProducer.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.NewString(),
		OrderUUID:       order.UUID,
		UserUUID:        order.User,
		PaymentMethod:   string(method),
		TransactionUUID: txID,
	}); err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err), zap.String("order_uuid", order.UUID))
	}
	return txID, nil
}
