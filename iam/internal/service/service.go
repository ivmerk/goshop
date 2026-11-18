package service

import (
	"context"

	"github.com/ivmerk/goshop/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, in *model.RegisterInput) (string, error)
	GetUser(ctx context.Context, uuid string) (*model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (*model.UserInfo, error)
}
