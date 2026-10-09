package v1

import (
	"github.com/ivmerk/goshop/payment/internal/service"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

var _ paymentV1.PaymentServiceServer = (*PaymentAPI)(nil)

type PaymentAPI struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewPaymentAPI(paymentService service.PaymentService) *PaymentAPI {
	return &PaymentAPI{
		paymentService: paymentService,
	}
}
