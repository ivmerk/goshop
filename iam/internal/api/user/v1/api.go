package v1

import (
	"github.com/ivmerk/goshop/iam/internal/service"
	userV1 "github.com/ivmerk/goshop/shared/pkg/proto/user/v1"
)

var _ userV1.UserServiceServer = (*api)(nil)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{userService: userService}
}
