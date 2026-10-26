package order

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
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

func TestCreate_EmptyUUID(t *testing.T) {
	repo, mock := newTestRepo(t)

	order := &model.Order{
		UUID:    "",
		User:    userUUID,
		Parts:   []string{partUUID1, partUUID2},
		Total:   15700.5,
		Payment: model.PaymentMethodCard,
		Status:  model.OrderStatusPendingPayment,
	}

	got, err := repo.Create(context.Background(), order)

	require.ErrorIs(t, err, model.ErrInvalidOrderUUID)
	require.Empty(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_AlreadyExists(t *testing.T) {
	repo, mock := newTestRepo(t)

	order := &model.Order{
		UUID:    newOrderUUID,
		User:    userUUID,
		Parts:   []string{partUUID1, partUUID2},
		Total:   15700.5,
		Payment: model.PaymentMethodCard,
		Status:  model.OrderStatusPendingPayment,
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	got, err := repo.Create(context.Background(), order)

	require.ErrorIs(t, err, model.ErrOrderAlreadyExists)
	require.Empty(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DBError(t *testing.T) {
	repo, mock := newTestRepo(t)
	dbErr := errors.New("connection lost")

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnError(dbErr)

	got, err := repo.Create(context.Background(), &model.Order{UUID: newOrderUUID})

	require.ErrorIs(t, err, dbErr)
	require.NotErrorIs(t, err, model.ErrOrderAlreadyExists)
	require.Empty(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNewOrderRepository(t *testing.T) {
	_, mock := newTestRepo(t)

	require.NotNil(t, NewOrderRepository(mock))
}
