package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	repoMocks "github.com/ivmerk/goshop/iam/internal/repository/mocks"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

const (
	userUUID    = "11111111-1111-1111-1111-111111111111"
	sessionUUID = "22222222-2222-2222-2222-222222222222"
	login       = "john"
	email       = "john@example.com"
	password    = "SecurePassword123!"
	sessionTTL  = 24 * time.Hour
)

type ServiceSuite struct {
	suite.Suite

	ctx         context.Context
	userRepo    *repoMocks.UserRepository
	sessionRepo *repoMocks.SessionRepository
	service     *service

	passwordHash string
}

func (s *ServiceSuite) SetupSuite() {
	logger.SetNopLogger()

	// MinCost, чтобы тесты не тратили время на хеширование
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	s.Require().NoError(err)
	s.passwordHash = string(hash)
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = repoMocks.NewUserRepository(s.T())
	s.sessionRepo = repoMocks.NewSessionRepository(s.T())
	s.service = NewService(s.userRepo, s.sessionRepo, sessionTTL)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
