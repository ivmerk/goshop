package app

import (
	"context"

	paymentAPI "github.com/ivmerk/goshop/payment/internal/api/payment/v1"
	"github.com/ivmerk/goshop/payment/internal/service"
	paymentSvc "github.com/ivmerk/goshop/payment/internal/service/payment"
	paymentV1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer
	paymentSvc   service.PaymentService
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) PaymentV1API(ctx context.Context) (paymentV1.PaymentServiceServer, error) {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentAPI.NewPaymentAPI(d.PaymentService(ctx))
	}
	return d.paymentV1API, nil
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentSvc == nil {
		d.paymentSvc = paymentSvc.NewService()
	}
	return d.paymentSvc
}
