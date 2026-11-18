package converter

import (
	"github.com/ivmerk/goshop/iam/internal/model"
	commonV1 "github.com/ivmerk/goshop/shared/pkg/proto/common/v1"
)

func UserInfoToProto(info *model.UserInfo) *commonV1.UserInfo {
	return &commonV1.UserInfo{
		UserUuid: info.UUID,
		Login:    info.Login,
		Email:    info.Email,
	}
}
