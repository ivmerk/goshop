package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ivmerk/goshop/iam/internal/model"
	repoConverter "github.com/ivmerk/goshop/iam/internal/repository/converter"
)

func (r *repo) Create(ctx context.Context, session *model.Session, ttl time.Duration) error {
	data, err := json.Marshal(repoConverter.ToRepositorySession(session))
	if err != nil {
		return err
	}

	return r.cache.SetWithTTL(ctx, sessionKey(session.UUID), data, ttl)
}
