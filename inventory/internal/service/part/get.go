package part

import (
	"context"

	"github.com/ivmerk/goshop/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.repository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return part, nil
}
