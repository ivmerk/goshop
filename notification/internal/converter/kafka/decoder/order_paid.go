package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/ivmerk/goshop/notification/internal/model"
	eventsV1 "github.com/ivmerk/goshop/shared/pkg/proto/events/v1"
)

func OrderPaidToModel(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, err
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.GetEventUuid(),
		OrderUUID:       pb.GetOrderUuid(),
		UserUUID:        pb.GetUserUuid(),
		PaymentMethod:   pb.GetPaymentMethod(),
		TransactionUUID: pb.GetTransactionUuid(),
	}, nil
}
