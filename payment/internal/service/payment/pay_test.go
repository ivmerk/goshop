package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/payment/internal/model"
)

const (
	validOrderUUID = "11111111-1111-1111-1111-111111111111"
	validUserUUID  = "22222222-2222-2222-2222-222222222222"
)

func TestPay(t *testing.T) {
	svc := NewService()
	tests := []struct {
		name      string
		orderUUID string
		userUUID  string
		method    model.PaymentMethod
		wantErr   error
	}{
		{name: "card", orderUUID: validOrderUUID, userUUID: validUserUUID, method: model.PaymentMethodCard},
		{name: "sbp", orderUUID: validOrderUUID, userUUID: validUserUUID, method: model.PaymentMethodSbp},
		{name: "credit_card", orderUUID: validOrderUUID, userUUID: validUserUUID, method: model.PaymentMethodCreditCard},
		{name: "investor_money", orderUUID: validOrderUUID, userUUID: validUserUUID, method: model.PaymentMethodInvestorMoney},
		{name: "invalid order uuid", orderUUID: "abc", userUUID: validUserUUID, method: model.PaymentMethodCard, wantErr: model.ErrInvalidOrderUUID},
		{name: "invalid user uuid", orderUUID: validOrderUUID, userUUID: "", method: model.PaymentMethodCard, wantErr: model.ErrInvalidUserUUID},
		{name: "unknown method", orderUUID: validOrderUUID, userUUID: validUserUUID, method: model.PaymentMethodUnknown, wantErr: model.ErrUnsupportedPaymentMethod},
		{name: "not listed payment method", orderUUID: validOrderUUID, userUUID: validUserUUID, method: "Bitcoin", wantErr: model.ErrUnsupportedPaymentMethod},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := svc.Pay(context.Background(), tc.orderUUID, tc.userUUID, tc.method)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Empty(t, id)
				return
			}
			require.NoError(t, err)
			_, perr := uuid.Parse(id)
			require.NoError(t, perr)
		})
	}
}

func TestPay_UniqueTransactions(t *testing.T) {
	svc := NewService()
	a, _ := svc.Pay(context.Background(), validOrderUUID, validUserUUID, model.PaymentMethodCard)
	b, _ := svc.Pay(context.Background(), validOrderUUID, validUserUUID, model.PaymentMethodCard)
	require.NotEqual(t, a, b)
}
