package order

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/order/internal/model"
)

const updateQuery = "UPDATE orders"

func TestUpdate_Success(t *testing.T) {
	repo, mock := newTestRepo(t)
	tx := transactionUUID

	mock.ExpectExec(updateQuery).
		WithArgs(pendingOrderUUID, &tx, "PAID", "SBP").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	got, err := repo.Update(context.Background(), pendingOrderUUID, &model.OrderUpdateInfo{

		Transaction: &tx,
		Status:      model.OrderStatusPaid,
		Payment:     model.PaymentMethodSbp,
	})

	require.NoError(t, err)
	require.Equal(t, pendingOrderUUID, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_EmptyFieldsArePassedAsKeepOld(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectExec(updateQuery).
		WithArgs(paidOrderUUID, (*string)(nil), "", "").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	got, err := repo.Update(context.Background(), paidOrderUUID, &model.OrderUpdateInfo{})

	require.NoError(t, err)
	require.Equal(t, paidOrderUUID, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectExec(updateQuery).
		WithArgs(pendingOrderUUID, (*string)(nil), "CANCELLED", "").
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	got, err := repo.Update(context.Background(), pendingOrderUUID, &model.OrderUpdateInfo{
		Status: model.OrderStatusCancelled,
	})

	require.ErrorIs(t, err, model.ErrOrderNotFound)
	require.Empty(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_DBError(t *testing.T) {
	repo, mock := newTestRepo(t)
	dbErr := errors.New("connection lost")

	mock.ExpectExec(updateQuery).
		WithArgs(pendingOrderUUID, (*string)(nil), "CANCELLED", "").
		WillReturnError(dbErr)

	got, err := repo.Update(context.Background(), pendingOrderUUID, &model.OrderUpdateInfo{
		Status: model.OrderStatusCancelled,
	})

	require.ErrorIs(t, err, dbErr)
	require.NotErrorIs(t, err, model.ErrOrderNotFound)
	require.Empty(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}
