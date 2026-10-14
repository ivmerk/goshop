package payment

import (
	"context"

	"github.com/ivmerk/goshop/payment/internal/model"

	"github.com/google/uuid"
)

func (s *service) Pay(_ context.Context, orderUUID string, userUUID string, paymentMethod model.PaymentMethod) (string, error) {

	_, err := uuid.Parse(userUUID)
	if err != nil {
		return "", model.ErrInvalidUserUUID
	}

	_, err = uuid.Parse(orderUUID)
	if err != nil {
		return "", model.ErrInvalidOrderUUID
	}

	if !s.isPaymentMethodSupported(paymentMethod) {
		return "", model.ErrUnsupportedPaymentMethod
	}
	transactionID := uuid.NewString()
	return transactionID, nil
}

func (s *service) isPaymentMethodSupported(paymentMethod model.PaymentMethod) bool {
	switch paymentMethod {
	case model.PaymentMethodCard, model.PaymentMethodSbp, model.PaymentMethodCreditCard, model.PaymentMethodInvestorMoney:
		return true
	default:
		return false
	}
}
