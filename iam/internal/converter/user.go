package converter

import (
	"strings"

	"github.com/ivmerk/goshop/iam/internal/model"
	commonV1 "github.com/ivmerk/goshop/shared/pkg/proto/common/v1"
	userV1 "github.com/ivmerk/goshop/shared/pkg/proto/user/v1"
)

func RegisterInputFromProto(req *userV1.RegisterRequest) *model.RegisterInput {
	return &model.RegisterInput{
		Login:               req.GetLogin(),
		Password:            req.GetPassword(),
		Email:               req.GetEmail(),
		NotificationMethods: notificationMethodsFromProto(req.GetNotificationMethods()),
	}
}

func UserToGetUserResponse(user *model.User) *userV1.GetUserResponse {
	return &userV1.GetUserResponse{
		UserUuid:            user.UUID,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethodsToProto(user.NotificationMethods),
	}
}

func notificationMethodsFromProto(in []*commonV1.NotificationMethod) []model.NotificationMethod {
	out := make([]model.NotificationMethod, 0, len(in))
	for _, m := range in {
		out = append(out, model.NotificationMethod{
			ProviderName: providerNameFromProto(m.GetProviderName()),
			Target:       m.GetTarget(),
		})
	}

	return out
}

func notificationMethodsToProto(in []model.NotificationMethod) []*commonV1.NotificationMethod {
	out := make([]*commonV1.NotificationMethod, 0, len(in))
	for _, m := range in {
		out = append(out, &commonV1.NotificationMethod{
			ProviderName: providerNameToProto(m.ProviderName),
			Target:       m.Target,
		})
	}

	return out
}

// В домене имя провайдера — строка в нижнем регистре (telegram, email, push).
func providerNameFromProto(p commonV1.ProviderName) string {
	if p == commonV1.ProviderName_PROVIDER_NAME_UNSPECIFIED {
		return ""
	}

	return strings.ToLower(p.String())
}

func providerNameToProto(name string) commonV1.ProviderName {
	if v, ok := commonV1.ProviderName_value[strings.ToUpper(name)]; ok {
		return commonV1.ProviderName(v)
	}

	return commonV1.ProviderName_PROVIDER_NAME_UNSPECIFIED
}
