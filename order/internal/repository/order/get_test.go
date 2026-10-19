package order

import (
	"context"
	"testing"

	"github.com/ivmerk/goshop/order/internal/model"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		uuid       string
		wantStatus model.OrderStatus
		wantTx     bool
		wantErr    error
	}{
		{name: "pending_order", uuid: pendingOrderUUID, wantStatus: model.OrderStatusPendingPayment},
		{name: "paid_order", uuid: paidOrderUUID, wantStatus: model.OrderStatusPaid, wantTx: true},
		{name: "not_found", uuid: "99999999-9999-9999-9999-999999999999", wantErr: model.ErrOrderNotFound},
		{name: "empty_uuid", uuid: "", wantErr: model.ErrOrderNotFound},
	}

	rep := newTestRepo()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := rep.Get(context.Background(), tc.uuid)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.uuid, got.UUID)
			require.Equal(t, userUUID, got.User)
			require.Equal(t, tc.wantStatus, got.Status)
			if tc.wantTx {
				require.NotNil(t, got.Transaction)
				require.Equal(t, transactionUUID, *got.Transaction)
			} else {
				require.Nil(t, got.Transaction)
			}
		})
	}
}

func TestGet_ReturnsAllFields(t *testing.T) {
	rep := newTestRepo()

	got, err := rep.Get(context.Background(), paidOrderUUID)

	require.NoError(t, err)
	require.Equal(t, []string{partUUID1}, got.Parts)
	require.InDelta(t, 12500.5, got.Total, 0.0001)
	require.Equal(t, model.PaymentMethodCard, got.Payment)
	require.Equal(t, model.OrderStatusPaid, got.Status)
}

func TestGet_ReturnsCopy(t *testing.T) {
	rep := newTestRepo()

	got, err := rep.Get(context.Background(), paidOrderUUID)
	require.NoError(t, err)

	got.Parts[0] = "changed"
	*got.Transaction = "changed"
	got.User = "changed"

	stored := rep.orders[paidOrderUUID]
	require.Equal(t, partUUID1, stored.Parts[0])
	require.Equal(t, transactionUUID, *stored.Transaction)
	require.Equal(t, userUUID, stored.User)
}
