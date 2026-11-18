package v1

import (
	"context"

	iamAPI "github.com/ivmerk/goshop/iam/internal/api"
	"github.com/ivmerk/goshop/iam/internal/converter"
	userV1 "github.com/ivmerk/goshop/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	uuid, err := a.userService.Register(ctx, converter.RegisterInputFromProto(req))
	if err != nil {
		return nil, iamAPI.ToGRPCError(err)
	}

	return &userV1.RegisterResponse{UserUuid: uuid}, nil
}
