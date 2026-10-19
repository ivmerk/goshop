package order

import (
	"time"

	repoModel "github.com/ivmerk/goshop/order/internal/repository/model"
)

const (
	pendingOrderUUID = "11111111-1111-1111-1111-111111111111"
	paidOrderUUID    = "22222222-2222-2222-2222-222222222222"
	userUUID         = "33333333-3333-3333-3333-333333333333"
	partUUID1        = "44444444-4444-4444-4444-444444444444"
	partUUID2        = "55555555-5555-5555-5555-555555555555"
	transactionUUID  = "66666666-6666-6666-6666-666666666666"
)

func newTestRepo() *repository {
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	tx := transactionUUID

	return &repository{orders: map[string]*repoModel.Order{
		pendingOrderUUID: {
			UUID:      pendingOrderUUID,
			User:      userUUID,
			Parts:     []string{partUUID1, partUUID2},
			Total:     15700.5,
			Status:    repoModel.OrderStatus("PENDING_PAYMENT"),
			CreatedAt: createdAt,
		},
		paidOrderUUID: {
			UUID:        paidOrderUUID,
			User:        userUUID,
			Parts:       []string{partUUID1},
			Total:       12500.5,
			Payment:     repoModel.PaymentMethod("CARD"),
			Transaction: &tx,
			Status:      repoModel.OrderStatus("PAID"),
			CreatedAt:   createdAt,
		},
	}}
}
