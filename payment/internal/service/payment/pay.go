package payment

import (
	"context"

	"github.com/ivmerk/goshop/payment/internal/model"

	"github.com/google/uuid"
)

func (s *service) Pay(_ context.Context, _ string, _ string, _ model.PaymentMethod) (string, error) {

	transactionID := uuid.NewString()
	return transactionID, nil
}
