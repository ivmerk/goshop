package v1

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/ivmerk/goshop/order/internal/model"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "part_uuids must not be empty",
		}, nil
	}

	partUUIDs := make([]string, 0, len(req.PartUuids))
	for _, id := range req.PartUuids {
		partUUIDs = append(partUUIDs, id.String())
	}

	prices, err := h.inventory.ListPartPrices(ctx, partUUIDs)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get parts from inventory",
		}, nil
	}

	var total float64
	for _, id := range partUUIDs {
		price, ok := prices[id]
		if !ok {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "part not found: " + id,
			}, nil
		}
		total += price
	}

	orderUUID := uuid.New()
	h.storage.AddOrder(&model.Order{
		UUID:   orderUUID.String(),
		User:   req.UserUUID.String(),
		Parts:  partUUIDs,
		Total:  total,
		Status: model.OrderStatusPendingPayment,
	})

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: total,
	}, nil
}
