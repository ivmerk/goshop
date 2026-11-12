package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/notification/internal/model"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

func TestOrderPaidToModel(t *testing.T) {
	data, err := proto.Marshal(&eventsV1.OrderPaid{
		EventUuid:       "event-1",
		OrderUuid:       "order-1",
		UserUuid:        "user-1",
		PaymentMethod:   "CARD",
		TransactionUuid: "tx-1",
	})
	require.NoError(t, err)

	got, err := OrderPaidToModel(data)

	require.NoError(t, err)
	require.Equal(t, model.OrderPaidEvent{
		EventUUID:       "event-1",
		OrderUUID:       "order-1",
		UserUUID:        "user-1",
		PaymentMethod:   "CARD",
		TransactionUUID: "tx-1",
	}, got)
}

func TestShipAssembledToModel(t *testing.T) {
	data, err := proto.Marshal(&eventsV1.ShipAssembled{
		EventUuid:    "event-1",
		OrderUuid:    "order-1",
		UserUuid:     "user-1",
		BuildTimeSec: 7,
	})
	require.NoError(t, err)

	got, err := ShipAssembledToModel(data)

	require.NoError(t, err)
	require.Equal(t, model.ShipAssembledEvent{
		EventUUID:    "event-1",
		OrderUUID:    "order-1",
		UserUUID:     "user-1",
		BuildTimeSec: 7,
	}, got)
}

func TestDecodeInvalidData(t *testing.T) {
	_, err := OrderPaidToModel([]byte("garbage"))
	require.Error(t, err)

	_, err = ShipAssembledToModel([]byte("garbage"))
	require.Error(t, err)
}
