package order

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	parts := []string{partUUID1, partUUID2}
	s.inventoryClient.On("ListPartPrices", s.ctx, parts).
		Return(map[string]float64{partUUID1: 100.5, partUUID2: 200}, nil)
	s.orderRepository.On("Create", s.ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.UUID != "" && o.User == userUUID && o.Total == 300.5 &&
			o.Status == model.OrderStatusPendingPayment
	})).Return("any", nil)

	order, err := s.service.Create(s.ctx, userUUID, parts)

	s.Require().NoError(err)
	s.Require().NotEmpty(order.UUID)
	s.Require().InDelta(300.5, order.Total, 0.0001)
	s.Require().Equal(model.OrderStatusPendingPayment, order.Status)
	s.Require().Equal(parts, order.Parts)
}

func (s *ServiceSuite) TestCreateEmptyParts() {
	order, err := s.service.Create(s.ctx, userUUID, nil)

	s.Require().ErrorIs(err, model.ErrEmptyParts)
	s.Require().Nil(order)
}

func (s *ServiceSuite) TestCreateInventoryError() {
	invErr := errors.New("inventory down")
	s.inventoryClient.On("ListPartPrices", s.ctx, []string{partUUID1}).Return(nil, invErr)

	order, err := s.service.Create(s.ctx, userUUID, []string{partUUID1})

	s.Require().ErrorIs(err, invErr)
	s.Require().Nil(order)
}

func (s *ServiceSuite) TestCreatePartNotFound() {
	s.inventoryClient.On("ListPartPrices", s.ctx, []string{partUUID1, partUUID2}).
		Return(map[string]float64{partUUID1: 100}, nil) // второй детали нет

	order, err := s.service.Create(s.ctx, userUUID, []string{partUUID1, partUUID2})

	s.Require().ErrorIs(err, model.ErrPartNotFound)
	s.Require().ErrorContains(err, partUUID2)
	s.Require().Nil(order)
}

func (s *ServiceSuite) TestCreateRepositoryError() {
	repoErr := errors.New("repo failure")
	s.inventoryClient.On("ListPartPrices", s.ctx, []string{partUUID1}).
		Return(map[string]float64{partUUID1: 100}, nil)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return("", repoErr)

	order, err := s.service.Create(s.ctx, userUUID, []string{partUUID1})

	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(order)
}
