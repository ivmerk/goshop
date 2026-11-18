package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ivmerk/goshop/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	if login == "" || password == "" {
		return "", fmt.Errorf("%w: login and password are required", model.ErrInvalidArgument)
	}

	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		// Не раскрываем, существует ли такой логин.
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", model.ErrInvalidCredentials
	}

	now := time.Now()
	session := &model.Session{
		UUID:      uuid.NewString(),
		UserUUID:  user.UUID,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}

	if err := s.sessionRepo.Create(ctx, session, s.sessionTTL); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return session.UUID, nil
}
