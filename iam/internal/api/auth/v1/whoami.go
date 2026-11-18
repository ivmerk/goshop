package v1

import (
	"context"

	iamAPI "github.com/ivmerk/goshop/iam/internal/api"
	"github.com/ivmerk/goshop/iam/internal/converter"
	authV1 "github.com/ivmerk/goshop/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	info, err := a.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, iamAPI.ToGRPCError(err)
	}

	return &authV1.WhoamiResponse{UserInfo: converter.UserInfoToProto(info)}, nil
}
