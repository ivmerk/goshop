package session

import (
	"context"
	"encoding/json"
	"errors"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/ivmerk/goshop/iam/internal/model"
	repoConverter "github.com/ivmerk/goshop/iam/internal/repository/converter"
	repoModel "github.com/ivmerk/goshop/iam/internal/repository/model"
)

func (r *repo) Get(ctx context.Context, uuid string) (*model.Session, error) {
	data, err := r.cache.Get(ctx, sessionKey(uuid))
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	var s repoModel.Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return repoConverter.FromRepositorySession(&s), nil
}
