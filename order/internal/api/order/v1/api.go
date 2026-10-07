package v1

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

type OrderHandler struct {
	storage   *model.OrderStorage
	payment   PaymentClient
	inventory InventoryClient
}

type InventoryClient interface {
	ListPartPrices(ctx context.Context, uuids []string) (map[string]float64, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (transactionUUID string, err error)
}

func NewOrderHandler(storage *model.OrderStorage, inventory InventoryClient, payment PaymentClient) *OrderHandler {
	return &OrderHandler{storage: storage, inventory: inventory, payment: payment}
}
