package service

import (
	"context"

	"github.com/ivmerk/goshop/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}
