package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/ivmerk/goshop/order/internal/model"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partUUIDs := make([]string, 0, len(req.PartUuids))
	for _, id := range req.PartUuids {
		partUUIDs = append(partUUIDs, id.String())
	}
	order, err := h.orderService.Create(ctx, req.UserUUID.String(), partUUIDs)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyParts), errors.Is(err, model.ErrPartNotFound):
			return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "failed to create order"}, nil
		}
	}
	orderUUID, err := uuid.Parse(order.UUID)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid order uuid"}, nil
	}
	return &orderV1.CreateOrderResponse{OrderUUID: orderUUID, TotalPrice: order.Total}, nil
}
