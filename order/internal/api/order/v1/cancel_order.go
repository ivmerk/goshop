package v1

import (
	"context"
	"net/http"

	"github.com/ivmerk/goshop/order/internal/models"
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
	case models.OrderStatusPaid:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "order is already paid and cannot be canceled",
		}, nil
	case models.OrderStatusPendingPayment:
		order.Status = models.OrderStatusCancelled
		h.storage.UpdateOrder(order)
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
