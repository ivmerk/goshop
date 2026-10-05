package v1

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) GetOrderById(_ context.Context, params orderV1.GetOrderByIdParams) (orderV1.GetOrderByIdRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID.String())
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	userUUID, err := uuid.Parse(order.User)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid user uuid"}, nil
	}

	partUUIDs := make([]uuid.UUID, 0, len(order.Parts))
	for _, p := range order.Parts {
		id, err := uuid.Parse(p)
		if err != nil {
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid part uuid"}, nil
		}
		partUUIDs = append(partUUIDs, id)
	}

	resp := &orderV1.GetOrderResponse{
		OrderUUID:  params.OrderUUID,
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: order.Total,
		Status:     orderV1.OrderStatus(order.Status),
	}

	// Эти поля есть только у оплаченного заказа.
	if order.Transaction != "" {
		txID, err := uuid.Parse(order.Transaction)
		if err != nil {
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid transaction uuid"}, nil
		}
		resp.TransactionUUID = orderV1.NewOptUUID(txID)
	}
	if order.Payment != "" {
		resp.PaymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(order.Payment))
	}

	return resp, nil
}
