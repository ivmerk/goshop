package user

import (
	"github.com/ivmerk/goshop/iam/internal/repository"
	def "github.com/ivmerk/goshop/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *service {
	return &service{userRepo: userRepo}
}
