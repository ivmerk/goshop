package session

import (
	"github.com/ivmerk/goshop/iam/internal/repository"
	"github.com/ivmerk/goshop/platform/pkg/cache"
)

var _ repository.SessionRepository = (*repo)(nil)

const keyPrefix = "session:"

type repo struct {
	cache cache.RedisClient
}

func NewSessionRepository(c cache.RedisClient) *repo {
	return &repo{cache: c}
}

func sessionKey(uuid string) string {
	return keyPrefix + uuid
}
