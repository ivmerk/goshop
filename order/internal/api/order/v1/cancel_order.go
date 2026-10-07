package v1

import (
	"context"
	"net/http"

	"github.com/ivmerk/goshop/order/internal/model"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID.String())
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "order is already paid and cannot be canceled",
		}, nil
	case model.OrderStatusPendingPayment:
		order.Status = model.OrderStatusCancelled
		h.storage.UpdateOrder(order)
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
