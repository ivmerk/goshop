package user

import (
	"context"
	"fmt"

	"github.com/ivmerk/goshop/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, uuid string) (*model.User, error) {
	if uuid == "" {
		return nil, fmt.Errorf("%w: user uuid is required", model.ErrInvalidArgument)
	}

	return s.userRepo.GetByUUID(ctx, uuid)
}
