package converter

import (
	"github.com/ivmerk/goshop/iam/internal/model"
	repoModel "github.com/ivmerk/goshop/iam/internal/repository/model"
)

func ToRepositoryUser(user *model.User) *repoModel.User {
	methods := make([]repoModel.NotificationMethod, 0, len(user.NotificationMethods))
	for _, m := range user.NotificationMethods {
		methods = append(methods, repoModel.NotificationMethod{ProviderName: m.ProviderName, Target: m.Target})
	}

	return &repoModel.User{
		UUID:                user.UUID,
		Login:               user.Login,
		Email:               user.Email,
		PasswordHash:        user.PasswordHash,
		NotificationMethods: methods,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func FromRepositoryUser(user *repoModel.User) *model.User {
	methods := make([]model.NotificationMethod, 0, len(user.NotificationMethods))
	for _, m := range user.NotificationMethods {
		methods = append(methods, model.NotificationMethod{ProviderName: m.ProviderName, Target: m.Target})
	}

	return &model.User{
		UUID:                user.UUID,
		Login:               user.Login,
		Email:               user.Email,
		PasswordHash:        user.PasswordHash,
		NotificationMethods: methods,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}
