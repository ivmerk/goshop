package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivmerk/goshop/payment/internal/model"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

func (p *PaymentAPI) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionID, err := p.paymentService.Pay(ctx, req.GetOrderUuid(), req.GetUserUuid(), model.PaymentMethod(req.GetPaymentMethod().String()))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderUUID),
			errors.Is(err, model.ErrInvalidUserUUID),
			errors.Is(err, model.ErrUnsupportedPaymentMethod):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionID,
	}, nil
}
