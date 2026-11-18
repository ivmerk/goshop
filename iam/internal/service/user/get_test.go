package user

import "github.com/ivmerk/goshop/iam/internal/model"

func (s *ServiceSuite) TestGetUserSuccess() {
	want := &model.User{UUID: userUUID, Login: login, Email: email}
	s.userRepo.On("GetByUUID", s.ctx, userUUID).Return(want, nil)

	got, err := s.service.GetUser(s.ctx, userUUID)

	s.Require().NoError(err)
	s.Require().Equal(want, got)
}

func (s *ServiceSuite) TestGetUserEmptyUUID() {
	got, err := s.service.GetUser(s.ctx, "")

	s.Require().ErrorIs(err, model.ErrInvalidArgument)
	s.Require().Nil(got)
}

func (s *ServiceSuite) TestGetUserNotFound() {
	s.userRepo.On("GetByUUID", s.ctx, userUUID).Return(nil, model.ErrUserNotFound)

	got, err := s.service.GetUser(s.ctx, userUUID)

	s.Require().ErrorIs(err, model.ErrUserNotFound)
	s.Require().Nil(got)
}
