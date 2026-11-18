package auth

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/ivmerk/goshop/iam/internal/model"
)

func (s *ServiceSuite) TestLoginSuccess() {
	s.userRepo.On("GetByLogin", s.ctx, login).
		Return(&model.User{UUID: userUUID, Login: login, PasswordHash: s.passwordHash}, nil)
	s.sessionRepo.On("Create", s.ctx, mock.MatchedBy(func(sess *model.Session) bool {
		return sess.UUID != "" && sess.UserUUID == userUUID &&
			sess.ExpiresAt.After(sess.CreatedAt)
	}), sessionTTL).Return(nil)

	got, err := s.service.Login(s.ctx, login, password)

	s.Require().NoError(err)
	s.Require().NotEmpty(got)
}

func (s *ServiceSuite) TestLoginEmptyArguments() {
	for name, args := range map[string][2]string{
		"empty login":    {"", password},
		"empty password": {login, ""},
	} {
		s.Run(name, func() {
			got, err := s.service.Login(s.ctx, args[0], args[1])

			s.Require().ErrorIs(err, model.ErrInvalidArgument)
			s.Require().Empty(got)
		})
	}
}

func (s *ServiceSuite) TestLoginUnknownUser() {
	s.userRepo.On("GetByLogin", s.ctx, login).Return(nil, model.ErrUserNotFound)

	got, err := s.service.Login(s.ctx, login, password)

	// неизвестный логин неотличим от неверного пароля
	s.Require().ErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestLoginWrongPassword() {
	s.userRepo.On("GetByLogin", s.ctx, login).
		Return(&model.User{UUID: userUUID, Login: login, PasswordHash: s.passwordHash}, nil)

	got, err := s.service.Login(s.ctx, login, "wrong-password")

	s.Require().ErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestLoginUserRepoError() {
	repoErr := errors.New("db down")
	s.userRepo.On("GetByLogin", s.ctx, login).Return(nil, repoErr)

	got, err := s.service.Login(s.ctx, login, password)

	s.Require().ErrorIs(err, repoErr)
	s.Require().NotErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestLoginSessionRepoError() {
	repoErr := errors.New("redis down")
	s.userRepo.On("GetByLogin", s.ctx, login).
		Return(&model.User{UUID: userUUID, Login: login, PasswordHash: s.passwordHash}, nil)
	s.sessionRepo.On("Create", s.ctx, mock.Anything, sessionTTL).Return(repoErr)

	got, err := s.service.Login(s.ctx, login, password)

	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(got)
}
