package grpc

import (
	"context"
	"fmt"

	"github.com/ivmerk/goshop/order/internal/model"
	payment_v1 "github.com/ivmerk/goshop/shared/pkg/proto/payment/v1"
)

type PaymentClient struct {
	client payment_v1.PaymentServiceClient
}

func NewPaymentClient(client payment_v1.PaymentServiceClient) *PaymentClient {
	return &PaymentClient{client: client}
}

func (c *PaymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error) {
	pm, ok := payment_v1.PaymentMethod_value[string(method)]
	if !ok {
		return "", fmt.Errorf("unsupported payment method: %s", method)
	}

	resp, err := c.client.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: payment_v1.PaymentMethod(pm),
	})
	if err != nil {
		return "", fmt.Errorf("payment PayOrder: %w", err)
	}

	return resp.GetTransactionUuid(), nil
}
