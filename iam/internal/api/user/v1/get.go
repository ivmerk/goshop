package v1

import (
	"context"

	iamAPI "github.com/ivmerk/goshop/iam/internal/api"
	"github.com/ivmerk/goshop/iam/internal/converter"
	userV1 "github.com/ivmerk/goshop/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		return nil, iamAPI.ToGRPCError(err)
	}

	return converter.UserToGetUserResponse(user), nil
}
