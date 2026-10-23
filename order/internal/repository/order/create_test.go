package order

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/order/internal/model"
	repoModel "github.com/ivmerk/goshop/order/internal/repository/model"
)

const newOrderUUID = "77777777-7777-7777-7777-777777777777"

func TestCreate_Success(t *testing.T) {
	repo, mock := newTestRepo(t)

	order := &model.Order{
		UUID:    pendingOrderUUID,
		User:    userUUID,
		Parts:   []string{partUUID1, partUUID2},
		Total:   15700.5,
		Payment: model.PaymentMethodCard,
		Status:  model.OrderStatusPendingPayment,
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(pendingOrderUUID, // $1 uuid          string
			userUUID,                                 // $2 user          string
			[]string{partUUID1, partUUID2},           // $3 parts         []string
			15700.5,                                  // $4 total         float64
			repoModel.PaymentMethod("CARD"),          // $5 payment       repoModel.PaymentMethod
			(*string)(nil),                           // $6 transaction   *string, nil
			repoModel.OrderStatus("PENDING_PAYMENT"), // $7 status
		).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	_, err := repo.Create(context.Background(), order)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
