package v1

import (
	"github.com/ivmerk/goshop/order/internal/model"
)

type OrderHandler struct {
	storage   *model.OrderStorage
	payment   PaymentClient
	inventory InventoryClient
}

func NewOrderHandler(storage *model.OrderStorage, inventory InventoryClient, payment PaymentClient) *OrderHandler {
	return &OrderHandler{storage: storage, inventory: inventory, payment: payment}
}
