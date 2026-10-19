package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivmerk/goshop/order/internal/model"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	txID, err := h.orderService.Pay(ctx, params.OrderUUID.String(), model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
		case errors.Is(err, model.ErrInvalidPaymentMethod), errors.Is(err, model.ErrOrderNotPayable):
			return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "payment failed"}, nil
		}
	}

	transactionUUID, err := uuid.Parse(txID)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid transaction id"}, nil
	}
	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}
