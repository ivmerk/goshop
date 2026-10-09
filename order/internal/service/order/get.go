package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Order, error) {
	return &model.Order{UUID: &uuid}, nil
}
