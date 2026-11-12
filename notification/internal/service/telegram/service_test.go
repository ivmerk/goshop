package telegram

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	clientMocks "github.com/ivmerk/goshop/notification/internal/client/http/mocks"
	"github.com/ivmerk/goshop/notification/internal/model"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

const chatID int64 = 12345

func newTestService(t *testing.T) (*service, *clientMocks.TelegramClient) {
	t.Helper()
	logger.SetNopLogger()

	c := clientMocks.NewTelegramClient(t)
	s, err := NewService(c, chatID)
	require.NoError(t, err)

	return s, c
}

func TestSendOrderPaidNotification(t *testing.T) {
	s, c := newTestService(t)
	ctx := context.Background()

	c.On("SendMessage", ctx, chatID, "💳 Заказ оплачен!\n\nЗаказ: order-1\nПользователь: user-1\nСпособ оплаты: CARD\nТранзакция: tx-1\n").
		Return(nil)

	err := s.SendOrderPaidNotification(ctx, model.OrderPaidEvent{
		OrderUUID:       "order-1",
		UserUUID:        "user-1",
		PaymentMethod:   "CARD",
		TransactionUUID: "tx-1",
	})

	require.NoError(t, err)
}

func TestSendShipAssembledNotification(t *testing.T) {
	s, c := newTestService(t)
	ctx := context.Background()

	c.On("SendMessage", ctx, chatID, "📦 Заказ собран!\n\nЗаказ: order-1\nПользователь: user-1\nВремя сборки: 7 с\n").
		Return(nil)

	err := s.SendShipAssembledNotification(ctx, model.ShipAssembledEvent{
		OrderUUID:    "order-1",
		UserUUID:     "user-1",
		BuildTimeSec: 7,
	})

	require.NoError(t, err)
}

func TestSendNotificationClientError(t *testing.T) {
	s, c := newTestService(t)
	ctx := context.Background()
	sendErr := errors.New("telegram down")

	c.On("SendMessage", ctx, chatID, "📦 Заказ собран!\n\nЗаказ: order-1\nПользователь: user-1\nВремя сборки: 7 с\n").
		Return(sendErr)

	err := s.SendShipAssembledNotification(ctx, model.ShipAssembledEvent{
		OrderUUID:    "order-1",
		UserUUID:     "user-1",
		BuildTimeSec: 7,
	})

	require.ErrorIs(t, err, sendErr)
}
