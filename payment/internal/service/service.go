package service

import (
	"context"

	"github.com/ivmerk/goshop/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, orderUUID string, userUUID string, method model.PaymentMethod) (string, error)
}
