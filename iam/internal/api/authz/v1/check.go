package v1

import (
	"context"
	"errors"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	iamAPI "github.com/ivmerk/goshop/iam/internal/api"
	"github.com/ivmerk/goshop/iam/internal/model"
)

const (
	// Envoy передаёт имена заголовков в нижнем регистре.
	sessionUUIDHeader = "x-session-uuid"

	// Заголовки, которые Envoy добавит к запросу перед отправкой в апстрим.
	userUUIDHeader  = "x-user-uuid"
	userLoginHeader = "x-user-login"
)

// Check проверяет сессию из заголовка X-Session-Uuid.
// Отказ из-за отсутствия или недействительности сессии — это штатный ответ Denied (HTTP 401),
// а не ошибка gRPC. Ошибку gRPC возвращаем только при сбое самого IAM: Envoy в этом случае
// закрывает запрос (fail closed), а не пропускает его без проверки.
func (a *api) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	sessionUUID := req.GetAttributes().GetRequest().GetHttp().GetHeaders()[sessionUUIDHeader]
	if sessionUUID == "" {
		return denied("authentication required"), nil
	}

	info, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return denied("session not found or expired"), nil
		}
		return nil, iamAPI.ToGRPCError(err)
	}

	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: []*corev3.HeaderValueOption{
					header(userUUIDHeader, info.UUID),
					header(userLoginHeader, info.Login),
				},
			},
		},
	}, nil
}

func header(key, value string) *corev3.HeaderValueOption {
	return &corev3.HeaderValueOption{
		// Перезаписываем заголовок: клиент не должен подделать user-uuid своим значением.
		AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
		Header:       &corev3.HeaderValue{Key: key, Value: value},
	}
}

func denied(message string) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.Unauthenticated), Message: message},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: typev3.StatusCode_Unauthorized},
				Body:   `{"code":401,"message":"` + message + `"}`,
				Headers: []*corev3.HeaderValueOption{
					{Header: &corev3.HeaderValue{Key: "content-type", Value: "application/json"}},
				},
			},
		},
	}
}
