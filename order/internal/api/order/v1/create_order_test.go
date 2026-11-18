package v1

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/order/internal/model"
	"github.com/ivmerk/goshop/order/internal/service/mocks"
	grpcAuth "github.com/ivmerk/goshop/platform/pkg/middleware/grpc"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
	commonV1 "github.com/ivmerk/goshop/shared/pkg/proto/common/v1"
)

const (
	sessionUserUUID = "22222222-2222-2222-2222-222222222222"
	orderUUID       = "11111111-1111-1111-1111-111111111111"
)

func ctxWithUser(userUUID string) context.Context {
	return context.WithValue(context.Background(), grpcAuth.GetUserContextKey(), &commonV1.UserInfo{UserUuid: userUUID})
}

func TestCreateOrder_UsesUserFromSession(t *testing.T) {
	partID := uuid.New()
	svc := mocks.NewOrderService(t)
	svc.On("Create", ctxWithUser(sessionUserUUID), sessionUserUUID, []string{partID.String()}).
		Return(&model.Order{UUID: orderUUID, Total: 100}, nil)

	res, err := NewOrderHandler(svc).CreateOrder(ctxWithUser(sessionUserUUID),
		&orderV1.CreateOrderRequest{PartUuids: []uuid.UUID{partID}})

	require.NoError(t, err)
	resp, ok := res.(*orderV1.CreateOrderResponse)
	require.True(t, ok)
	require.Equal(t, orderUUID, resp.OrderUUID.String())
}

func TestCreateOrder_NoUserInContext(t *testing.T) {
	svc := mocks.NewOrderService(t) // сервис не должен вызываться

	res, err := NewOrderHandler(svc).CreateOrder(context.Background(),
		&orderV1.CreateOrderRequest{PartUuids: []uuid.UUID{uuid.New()}})

	require.NoError(t, err)
	_, ok := res.(*orderV1.InternalServerError)
	require.True(t, ok)
}
