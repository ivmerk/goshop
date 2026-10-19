package order

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (s *ServiceSuite) pendingOrder() *model.Order {
	return &model.Order{
		UUID:   orderUUID,
		User:   userUUID,
		Status: model.OrderStatusPendingPayment,
	}
}

func (s *ServiceSuite) TestPaySuccess() {
	tx := txUUID
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(s.pendingOrder(), nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, model.PaymentMethodCard).
		Return(txUUID, nil)
	s.orderRepository.On("Update", s.ctx, orderUUID, &model.OrderUpdateInfo{
		Transaction: &tx,
		Status:      model.OrderStatusPaid,
		Payment:     model.PaymentMethodCard,
	}).Return(orderUUID, nil)

	got, err := s.service.Pay(s.ctx, orderUUID, model.PaymentMethodCard)

	s.Require().NoError(err)
	s.Require().Equal(txUUID, got)
}

func (s *ServiceSuite) TestPayUnknownMethod() {
	got, err := s.service.Pay(s.ctx, orderUUID, model.PaymentMethodUnknown)

	s.Require().ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	got, err := s.service.Pay(s.ctx, orderUUID, model.PaymentMethodCard)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestPayOrderNotPayable() {
	for _, status := range []model.OrderStatus{model.OrderStatusPaid, model.OrderStatusCancelled} {
		s.Run(string(status), func() {
			s.orderRepository.On("Get", s.ctx, orderUUID).Unset()
			s.orderRepository.On("Get", s.ctx, orderUUID).
				Return(&model.Order{UUID: orderUUID, User: userUUID, Status: status}, nil).Once()

			got, err := s.service.Pay(s.ctx, orderUUID, model.PaymentMethodCard)

			s.Require().ErrorIs(err, model.ErrOrderNotPayable)
			s.Require().Empty(got)
		})
	}
}

func (s *ServiceSuite) TestPayPaymentError() {
	payErr := errors.New("payment down")
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(s.pendingOrder(), nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, model.PaymentMethodCard).
		Return("", payErr)

	got, err := s.service.Pay(s.ctx, orderUUID, model.PaymentMethodCard)

	s.Require().ErrorIs(err, payErr)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestPayUpdateError() {
	updateErr := errors.New("update failed")
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(s.pendingOrder(), nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, model.PaymentMethodSbp).
		Return(txUUID, nil)
	s.orderRepository.On("Update", s.ctx, orderUUID, mock.Anything).Return("", updateErr)

	_, err := s.service.Pay(s.ctx, orderUUID, model.PaymentMethodSbp)

	s.Require().ErrorIs(err, updateErr)
}
