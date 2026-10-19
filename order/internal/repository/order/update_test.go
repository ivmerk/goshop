package order

import (
	"context"
	"testing"
	"time"

	"github.com/ivmerk/goshop/order/internal/model"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	tx := "88888888-8888-8888-8888-888888888888"

	tests := []struct {
		name        string
		uuid        string
		info        *model.OrderUpdateInfo
		wantStatus  string
		wantPayment string
		wantTx      string
		wantErr     error
	}{
		{
			name: "pay_pending_order",
			uuid: pendingOrderUUID,
			info: &model.OrderUpdateInfo{
				Transaction: &tx,
				Status:      model.OrderStatusPaid,
				Payment:     model.PaymentMethodSbp,
			},
			wantStatus:  "PAID",
			wantPayment: "SBP",
			wantTx:      tx,
		},
		{
			name:       "cancel_pending_order",
			uuid:       pendingOrderUUID,
			info:       &model.OrderUpdateInfo{Status: model.OrderStatusCancelled},
			wantStatus: "CANCELLED",
		},
		{
			name:        "empty_status_keeps_old_one",
			uuid:        paidOrderUUID,
			info:        &model.OrderUpdateInfo{},
			wantStatus:  "PAID",
			wantPayment: "CARD",
			wantTx:      transactionUUID,
		},
		{
			name:    "order_not_found",
			uuid:    "99999999-9999-9999-9999-999999999999",
			info:    &model.OrderUpdateInfo{Status: model.OrderStatusPaid},
			wantErr: model.ErrOrderNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rep := newTestRepo()
			start := time.Now()

			got, err := rep.Update(context.Background(), tc.uuid, tc.info)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Empty(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.uuid, got)

			stored := rep.orders[tc.uuid]
			require.EqualValues(t, tc.wantStatus, stored.Status)
			require.EqualValues(t, tc.wantPayment, stored.Payment)
			if tc.wantTx == "" {
				require.Nil(t, stored.Transaction)
			} else {
				require.NotNil(t, stored.Transaction)
				require.Equal(t, tc.wantTx, *stored.Transaction)
			}
			require.NotNil(t, stored.UpdatedAt)
			require.False(t, stored.UpdatedAt.Before(start))
		})
	}
}

func TestUpdate_KeepsOtherFields(t *testing.T) {
	rep := newTestRepo()
	before := *rep.orders[paidOrderUUID]

	_, err := rep.Update(context.Background(), paidOrderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusCancelled})
	require.NoError(t, err)

	after := rep.orders[paidOrderUUID]
	require.Equal(t, before.User, after.User)
	require.Equal(t, before.Parts, after.Parts)
	require.Equal(t, before.Total, after.Total)
	require.Equal(t, before.CreatedAt, after.CreatedAt)
}

func TestUpdate_StoresCopyOfTransaction(t *testing.T) {
	rep := newTestRepo()
	tx := "88888888-8888-8888-8888-888888888888"

	_, err := rep.Update(context.Background(), pendingOrderUUID,
		&model.OrderUpdateInfo{Transaction: &tx})
	require.NoError(t, err)

	tx = "changed" // меняем исходную переменную после Update
	require.Equal(t, "88888888-8888-8888-8888-888888888888", *rep.orders[pendingOrderUUID].Transaction)
}
