package v1

import (
	"github.com/ivmerk/goshop/iam/internal/service"
	authV1 "github.com/ivmerk/goshop/shared/pkg/proto/auth/v1"
)

var _ authV1.AuthServiceServer = (*api)(nil)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{authService: authService}
}
