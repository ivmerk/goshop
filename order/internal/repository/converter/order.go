package converter

import (
	"slices"

	"github.com/ivmerk/goshop/order/internal/model"
	repositoryModel "github.com/ivmerk/goshop/order/internal/repository/model"
)

func ToRepositoryOrder(order *model.Order) *repositoryModel.Order {
	return &repositoryModel.Order{
		UUID:        order.UUID,
		User:        order.User,
		Parts:       slices.Clone(order.Parts),
		Total:       order.Total,
		Transaction: clonePtr(order.Transaction),
		Status:      repositoryModel.OrderStatus(order.Status),
		Payment:     repositoryModel.PaymentMethod(order.Payment),
	}
}

func FromRepositoryOrder(order *repositoryModel.Order) *model.Order {
	return &model.Order{
		UUID:        order.UUID,
		User:        order.User,
		Parts:       slices.Clone(order.Parts),
		Total:       order.Total,
		Transaction: clonePtr(order.Transaction),
		Status:      model.OrderStatus(order.Status),
		Payment:     model.PaymentMethod(order.Payment),
	}
}

func clonePtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}
