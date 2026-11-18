package repository

import (
	"context"
	"time"

	"github.com/ivmerk/goshop/iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (string, error)
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session, ttl time.Duration) error
	Get(ctx context.Context, uuid string) (*model.Session, error)
}
