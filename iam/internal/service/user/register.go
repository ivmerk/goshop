package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ivmerk/goshop/iam/internal/model"
)

// bcrypt принимает пароли не длиннее 72 байт.
const maxPasswordLen = 72

func (s *service) Register(ctx context.Context, in *model.RegisterInput) (string, error) {
	if in.Login == "" || in.Email == "" || in.Password == "" {
		return "", fmt.Errorf("%w: login, email and password are required", model.ErrInvalidArgument)
	}
	if len(in.Password) > maxPasswordLen {
		return "", fmt.Errorf("%w: password is too long", model.ErrInvalidArgument)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return s.userRepo.Create(ctx, &model.User{
		Login:               in.Login,
		Email:               in.Email,
		PasswordHash:        string(hash),
		NotificationMethods: in.NotificationMethods,
	})
}
