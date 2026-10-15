package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/ivmerk/goshop/order/internal/model"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if err := h.orderService.Cancel(ctx, params.OrderUUID.String()); err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{Code: http.StatusConflict, Message: "order is already paid and cannot be canceled"}, nil
		default:
			return nil, err
		}
	}
	return &orderV1.CancelOrderNoContent{}, nil

}
