package user

import (
	"errors"
	"strings"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/ivmerk/goshop/iam/internal/model"
)

func (s *ServiceSuite) TestRegisterSuccess() {
	methods := []model.NotificationMethod{{ProviderName: "telegram", Target: "12345"}}

	s.userRepo.On("Create", s.ctx, mock.MatchedBy(func(u *model.User) bool {
		// пароль хранится только в виде bcrypt-хеша
		return u.Login == login && u.Email == email &&
			u.PasswordHash != password &&
			bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil &&
			len(u.NotificationMethods) == 1
	})).Return(userUUID, nil)

	got, err := s.service.Register(s.ctx, &model.RegisterInput{
		Login: login, Password: password, Email: email, NotificationMethods: methods,
	})

	s.Require().NoError(err)
	s.Require().Equal(userUUID, got)
}

func (s *ServiceSuite) TestRegisterInvalidArgument() {
	tests := map[string]*model.RegisterInput{
		"empty login":    {Login: "", Password: password, Email: email},
		"empty password": {Login: login, Password: "", Email: email},
		"empty email":    {Login: login, Password: password, Email: ""},
		"too long":       {Login: login, Password: strings.Repeat("a", 73), Email: email},
	}

	for name, in := range tests {
		s.Run(name, func() {
			got, err := s.service.Register(s.ctx, in)

			s.Require().ErrorIs(err, model.ErrInvalidArgument)
			s.Require().Empty(got)
		})
	}
}

func (s *ServiceSuite) TestRegisterAlreadyExists() {
	s.userRepo.On("Create", s.ctx, mock.Anything).Return("", model.ErrUserAlreadyExists)

	got, err := s.service.Register(s.ctx, &model.RegisterInput{Login: login, Password: password, Email: email})

	s.Require().ErrorIs(err, model.ErrUserAlreadyExists)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestRegisterRepoError() {
	repoErr := errors.New("db down")
	s.userRepo.On("Create", s.ctx, mock.Anything).Return("", repoErr)

	got, err := s.service.Register(s.ctx, &model.RegisterInput{Login: login, Password: password, Email: email})

	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(got)
}
