package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/order/internal/model"
)

var getColumns = []string{"uuid", "user", "parts", "total", "payment", "transaction", "status", "created_at", "updated_at"}

func TestGet_Success(t *testing.T) {
	repo, mock := newTestRepo(t)
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	tx := transactionUUID

	mock.ExpectQuery("SELECT .* FROM orders").
		WithArgs(paidOrderUUID).
		WillReturnRows(pgxmock.NewRows(getColumns).AddRow(
			paidOrderUUID, userUUID, []string{partUUID1}, 12500.5, "CARD", &tx, "PAID", createdAt, &updatedAt,
		))

	got, err := repo.Get(context.Background(), paidOrderUUID)

	require.NoError(t, err)
	require.Equal(t, paidOrderUUID, got.UUID)
	require.Equal(t, userUUID, got.User)
	require.Equal(t, []string{partUUID1}, got.Parts)
	require.InDelta(t, 12500.5, got.Total, 0.0001)
	require.Equal(t, model.PaymentMethodCard, got.Payment)
	require.Equal(t, model.OrderStatusPaid, got.Status)
	require.NotNil(t, got.Transaction)
	require.Equal(t, transactionUUID, *got.Transaction)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_NullTransaction(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectQuery("SELECT .* FROM orders").
		WithArgs(pendingOrderUUID).
		WillReturnRows(pgxmock.NewRows(getColumns).AddRow(
			pendingOrderUUID, userUUID, []string{partUUID1, partUUID2}, 15700.5, "CARD", (*string)(nil),
			"PENDING_PAYMENT", time.Now(), (*time.Time)(nil),
		))

	got, err := repo.Get(context.Background(), pendingOrderUUID)

	require.NoError(t, err)
	require.Nil(t, got.Transaction)
	require.Equal(t, model.OrderStatusPendingPayment, got.Status)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_NotFound(t *testing.T) {
	repo, mock := newTestRepo(t)
	mock.ExpectQuery("SELECT .* FROM orders").WithArgs(pendingOrderUUID).WillReturnError(pgx.ErrNoRows)

	got, err := repo.Get(context.Background(), pendingOrderUUID)

	require.ErrorIs(t, err, model.ErrOrderNotFound)
	require.Nil(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_DBError(t *testing.T) {
	repo, mock := newTestRepo(t)
	dbErr := errors.New("connection lost")
	mock.ExpectQuery("SELECT .* FROM orders").WithArgs(pendingOrderUUID).WillReturnError(dbErr)

	got, err := repo.Get(context.Background(), pendingOrderUUID)

	require.ErrorIs(t, err, dbErr)
	require.NotErrorIs(t, err, model.ErrOrderNotFound)
	require.Nil(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}
