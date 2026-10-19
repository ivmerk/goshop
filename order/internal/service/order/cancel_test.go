package order

import (
	"errors"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	s.orderRepository.On("Get", s.ctx, orderUUID).
		Return(&model.Order{UUID: orderUUID, Status: model.OrderStatusPendingPayment}, nil)
	s.orderRepository.On("Update", s.ctx, orderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusCancelled}).
		Return(orderUUID, nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelAlreadyPaid() {
	s.orderRepository.On("Get", s.ctx, orderUUID).
		Return(&model.Order{UUID: orderUUID, Status: model.OrderStatusPaid}, nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelAlreadyCancelled() {
	s.orderRepository.On("Get", s.ctx, orderUUID).
		Return(&model.Order{UUID: orderUUID, Status: model.OrderStatusCancelled}, nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().NoError(err) // отмена идемпотентна
}

func (s *ServiceSuite) TestCancelUpdateError() {
	updateErr := errors.New("update failed")
	s.orderRepository.On("Get", s.ctx, orderUUID).
		Return(&model.Order{UUID: orderUUID, Status: model.OrderStatusPendingPayment}, nil)
	s.orderRepository.On("Update", s.ctx, orderUUID,
		&model.OrderUpdateInfo{Status: model.OrderStatusCancelled}).
		Return("", updateErr)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, updateErr)
}
