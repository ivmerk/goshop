package payment

import (
	"context"

	"github.com/ivmerk/goshop/payment/internal/model"

	"github.com/google/uuid"
)

func (s *service) Pay(_ context.Context, order_uuid string, _ string, paymentMethod model.PaymentMethod) (string, error) {
	_, err := uuid.Parse(order_uuid)
	if err != nil {
		return "", err
	}

	if paymentMethod == model.PaymentMethodUnknown {
		return "", model.ErrPaymentMethodUnknown
	}
	transactionID := uuid.NewString()
	return transactionID, nil
}
