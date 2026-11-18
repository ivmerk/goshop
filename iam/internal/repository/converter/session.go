package converter

import (
	"github.com/ivmerk/goshop/iam/internal/model"
	repoModel "github.com/ivmerk/goshop/iam/internal/repository/model"
)

func ToRepositorySession(s *model.Session) *repoModel.Session {
	return &repoModel.Session{
		UUID:      s.UUID,
		UserUUID:  s.UserUUID,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}

func FromRepositorySession(s *repoModel.Session) *model.Session {
	return &model.Session{
		UUID:      s.UUID,
		UserUUID:  s.UserUUID,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}
