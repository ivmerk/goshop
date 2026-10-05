package v1

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/models"
)

type OrderHandler struct {
	storage   *models.OrderStorage
	payment   PaymentClient
	inventory InventoryClient
}

type InventoryClient interface {
	ListPartPrices(ctx context.Context, uuids []string) (map[string]float64, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, method models.PaymentMethod) (transactionUUID string, err error)
}

func NewOrderHandler(storage *models.OrderStorage, inventory InventoryClient, payment PaymentClient) *OrderHandler {
	return &OrderHandler{storage: storage, inventory: inventory, payment: payment}
}
