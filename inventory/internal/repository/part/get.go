package part

import (
	"context"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoConverter "github.com/ivmerk/goshop/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, model.ErrPartNotFound
	}

	return repoConverter.FromRepositoryPart(part), nil
}
