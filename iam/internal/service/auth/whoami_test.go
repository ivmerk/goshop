package auth

import (
	"errors"

	"github.com/ivmerk/goshop/iam/internal/model"
)

func (s *ServiceSuite) TestWhoamiSuccess() {
	s.sessionRepo.On("Get", s.ctx, sessionUUID).Return(&model.Session{UUID: sessionUUID, UserUUID: userUUID}, nil)
	s.userRepo.On("GetByUUID", s.ctx, userUUID).
		Return(&model.User{UUID: userUUID, Login: login, Email: email, PasswordHash: s.passwordHash}, nil)

	got, err := s.service.Whoami(s.ctx, sessionUUID)

	s.Require().NoError(err)
	s.Require().Equal(&model.UserInfo{UUID: userUUID, Login: login, Email: email}, got)
}

func (s *ServiceSuite) TestWhoamiEmptySession() {
	got, err := s.service.Whoami(s.ctx, "")

	s.Require().ErrorIs(err, model.ErrInvalidArgument)
	s.Require().Nil(got)
}

func (s *ServiceSuite) TestWhoamiSessionNotFound() {
	s.sessionRepo.On("Get", s.ctx, sessionUUID).Return(nil, model.ErrSessionNotFound)

	got, err := s.service.Whoami(s.ctx, sessionUUID)

	s.Require().ErrorIs(err, model.ErrSessionNotFound)
	s.Require().Nil(got)
}

func (s *ServiceSuite) TestWhoamiUserError() {
	userErr := errors.New("db down")
	s.sessionRepo.On("Get", s.ctx, sessionUUID).Return(&model.Session{UUID: sessionUUID, UserUUID: userUUID}, nil)
	s.userRepo.On("GetByUUID", s.ctx, userUUID).Return(nil, userErr)

	got, err := s.service.Whoami(s.ctx, sessionUUID)

	s.Require().ErrorIs(err, userErr)
	s.Require().Nil(got)
}
