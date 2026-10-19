package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/order/internal/model"
)

const newOrderUUID = "77777777-7777-7777-7777-777777777777"

func TestNewOrderRepository(t *testing.T) {
	rep := NewOrderRepository()
	require.NotNil(t, rep)
	require.Empty(t, rep.orders)
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		order   *model.Order
		wantErr error
	}{
		{
			name: "success",
			order: &model.Order{
				UUID:   newOrderUUID,
				User:   userUUID,
				Parts:  []string{partUUID1},
				Total:  12500.5,
				Status: model.OrderStatusPendingPayment,
			},
		},
		{
			name:    "empty_uuid",
			order:   &model.Order{User: userUUID, Parts: []string{partUUID1}},
			wantErr: model.ErrInvalidOrderUUID,
		},
		{
			name:    "duplicate_uuid",
			order:   &model.Order{UUID: pendingOrderUUID, User: userUUID, Parts: []string{partUUID1}},
			wantErr: model.ErrOrderAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rep := newTestRepo()
			before := len(rep.orders)

			got, err := rep.Create(context.Background(), tc.order)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Empty(t, got)
				require.Len(t, rep.orders, before) // ничего не добавилось
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.order.UUID, got)
			require.Len(t, rep.orders, before+1)
		})
	}
}

func TestCreate_StoresOrderWithCreatedAt(t *testing.T) {
	rep := newTestRepo()
	start := time.Now()

	_, err := rep.Create(context.Background(), &model.Order{
		UUID:    newOrderUUID,
		User:    userUUID,
		Parts:   []string{partUUID1, partUUID2},
		Total:   100,
		Status:  model.OrderStatusPendingPayment,
		Payment: model.PaymentMethodCard,
	})
	require.NoError(t, err)

	stored := rep.orders[newOrderUUID]
	require.NotNil(t, stored)
	require.Equal(t, userUUID, stored.User)
	require.Equal(t, []string{partUUID1, partUUID2}, stored.Parts)
	require.InDelta(t, 100, stored.Total, 0.0001)
	require.EqualValues(t, "PENDING_PAYMENT", stored.Status)
	require.EqualValues(t, "CARD", stored.Payment)
	require.False(t, stored.CreatedAt.Before(start)) // выставляется репозиторием
}

func TestCreate_DuplicateKeepsOriginal(t *testing.T) {
	rep := newTestRepo()
	original := rep.orders[pendingOrderUUID]

	_, err := rep.Create(context.Background(), &model.Order{
		UUID: pendingOrderUUID, User: "someone-else", Total: 1,
	})
	require.ErrorIs(t, err, model.ErrOrderAlreadyExists)

	require.Same(t, original, rep.orders[pendingOrderUUID]) // тот же объект
	require.Equal(t, userUUID, rep.orders[pendingOrderUUID].User)
}

func TestCreate_StoresCopy(t *testing.T) {
	rep := newTestRepo()
	parts := []string{partUUID1}
	order := &model.Order{UUID: newOrderUUID, User: userUUID, Parts: parts}

	_, err := rep.Create(context.Background(), order)
	require.NoError(t, err)

	parts[0] = "changed" // меняем исходный срез после Create
	require.Equal(t, partUUID1, rep.orders[newOrderUUID].Parts[0])
}
