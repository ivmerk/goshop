package service

import (
	"context"
)

type PaymentService interface {
	Pay(ctx context.Context, uuid string, transactionID string) (string, error)
}
