// Package v1 реализует gRPC-контракт Envoy ext_authz (envoy.service.auth.v3.Authorization).
// Envoy перед проксированием защищённого запроса вызывает Check и пропускает его
// только если IAM подтвердил сессию.
package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/ivmerk/goshop/iam/internal/service"
)

var _ authv3.AuthorizationServer = (*api)(nil)

type api struct {
	authv3.UnimplementedAuthorizationServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{authService: authService}
}
