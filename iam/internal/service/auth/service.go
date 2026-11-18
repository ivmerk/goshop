package auth

import (
	"time"

	"github.com/ivmerk/goshop/iam/internal/repository"
	def "github.com/ivmerk/goshop/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	sessionTTL  time.Duration
}

func NewService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	sessionTTL time.Duration,
) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}
