package v1

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivmerk/goshop/order/internal/model"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID.String())
	if !ok {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}

	if req.PaymentMethod == orderV1.PaymentMethodUNKNOWN {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: "unknown payment method"}, nil
	}

	method := model.PaymentMethod(req.PaymentMethod)
	transactionUUID, err := h.payment.PayOrder(ctx, order.UUID, order.User, method)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "payment failed"}, nil
	}

	order.Status = model.OrderStatusPaid
	order.Transaction = transactionUUID
	order.Payment = method
	h.storage.UpdateOrder(order)

	txID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid transaction id"}, nil
	}
	return &orderV1.PayOrderResponse{TransactionUUID: txID}, nil
}
