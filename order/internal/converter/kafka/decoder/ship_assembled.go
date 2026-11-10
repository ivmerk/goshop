package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/order/internal/model"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

func ShipAssembledToModel(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, err
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
