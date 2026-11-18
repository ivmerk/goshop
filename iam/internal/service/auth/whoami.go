package auth

import (
	"context"
	"fmt"

	"github.com/ivmerk/goshop/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (*model.UserInfo, error) {
	if sessionUUID == "" {
		return nil, fmt.Errorf("%w: session uuid is required", model.ErrInvalidArgument)
	}

	session, err := s.sessionRepo.Get(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return nil, err
	}

	return &model.UserInfo{
		UUID:  user.UUID,
		Login: user.Login,
		Email: user.Email,
	}, nil
}
