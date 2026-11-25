package v1

import (
	"context"
	"errors"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivmerk/goshop/iam/internal/model"
	"github.com/ivmerk/goshop/iam/internal/service/mocks"
)

func checkRequest(headers map[string]string) *authv3.CheckRequest {
	return &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{Headers: headers},
			},
		},
	}
}

func TestCheck_MissingSession(t *testing.T) {
	svc := mocks.NewAuthService(t) // Whoami не должен вызываться
	resp, err := NewAPI(svc).Check(context.Background(), checkRequest(map[string]string{}))

	require.NoError(t, err)
	require.EqualValues(t, codes.Unauthenticated, resp.GetStatus().GetCode())
	require.Equal(t, typev3.StatusCode_Unauthorized, resp.GetDeniedResponse().GetStatus().GetCode())
}

func TestCheck_UnknownSession(t *testing.T) {
	svc := mocks.NewAuthService(t)
	svc.On("Whoami", context.Background(), "bad").Return(nil, model.ErrSessionNotFound)

	resp, err := NewAPI(svc).Check(context.Background(), checkRequest(map[string]string{"x-session-uuid": "bad"}))

	require.NoError(t, err)
	require.EqualValues(t, codes.Unauthenticated, resp.GetStatus().GetCode())
	require.Equal(t, typev3.StatusCode_Unauthorized, resp.GetDeniedResponse().GetStatus().GetCode())
}

func TestCheck_InternalErrorFailsClosed(t *testing.T) {
	svc := mocks.NewAuthService(t)
	svc.On("Whoami", context.Background(), "s1").Return(nil, errors.New("redis down"))

	resp, err := NewAPI(svc).Check(context.Background(), checkRequest(map[string]string{"x-session-uuid": "s1"}))

	require.Nil(t, resp)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestCheck_ValidSession(t *testing.T) {
	svc := mocks.NewAuthService(t)
	svc.On("Whoami", context.Background(), "s1").Return(&model.UserInfo{UUID: "u1", Login: "ivan"}, nil)

	resp, err := NewAPI(svc).Check(context.Background(), checkRequest(map[string]string{"x-session-uuid": "s1"}))

	require.NoError(t, err)
	require.EqualValues(t, codes.OK, resp.GetStatus().GetCode())

	got := map[string]string{}
	for _, h := range resp.GetOkResponse().GetHeaders() {
		got[h.GetHeader().GetKey()] = h.GetHeader().GetValue()
	}
	require.Equal(t, map[string]string{"x-user-uuid": "u1", "x-user-login": "ivan"}, got)
}
