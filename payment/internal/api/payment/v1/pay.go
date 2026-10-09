package v1

import (
	"context"

	"github.com/ivmerk/goshop/payment/internal/model"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

func (p *PaymentAPI) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionID, err := p.paymentService.Pay(ctx, req.GetOrderUuid(), req.GetUserUuid(), model.PaymentMethod(req.GetPaymentMethod().String()))
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionID,
	}, nil
}
