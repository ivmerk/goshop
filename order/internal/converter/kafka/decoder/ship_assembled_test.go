package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/order/internal/model"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

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

func TestShipAssembledToModelInvalidData(t *testing.T) {
	_, err := ShipAssembledToModel([]byte("garbage"))

	require.Error(t, err)
}
