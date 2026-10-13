package part

import (
	"context"

	"github.com/ivmerk/goshop/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	list, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return list, nil
}
