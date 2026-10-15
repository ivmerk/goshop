package grpc

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

type InventoryClient interface {
	ListPartPrices(ctx context.Context, uuids []string) (map[string]float64, error)
}
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}
