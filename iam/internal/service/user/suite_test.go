package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/ivmerk/goshop/iam/internal/repository/mocks"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

const (
	userUUID = "11111111-1111-1111-1111-111111111111"
	login    = "john"
	email    = "john@example.com"
	password = "SecurePassword123!"
)

type ServiceSuite struct {
	suite.Suite

	ctx      context.Context
	userRepo *repoMocks.UserRepository
	service  *service
}

func (s *ServiceSuite) SetupSuite() {
	logger.SetNopLogger()
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = repoMocks.NewUserRepository(s.T())
	s.service = NewService(s.userRepo)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
