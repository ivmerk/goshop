package converter

import (
	"github.com/ivmerk/goshop/order/internal/model"
	repositoryModel "github.com/ivmerk/goshop/order/internal/repository/model"
	"golang.org/x/exp/slices"
)

func ToRepositoryOrder(order *model.Order) *repositoryModel.Order {
	return &repositoryModel.Order{
		UUID:        order.UUID,
		User:        order.User,
		Parts:       slices.Clone(order.Parts),
		Total:       order.Total,
		Transaction: stringPtr(order.Transaction),
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
		Transaction: stringVal(order.Transaction),
		Status:      model.OrderStatus(order.Status),
		Payment:     model.PaymentMethod(order.Payment),
	}
}
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
