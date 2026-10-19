package order

import "github.com/ivmerk/goshop/order/internal/model"

func (s *ServiceSuite) TestGetSuccess() {
	want := &model.Order{UUID: orderUUID, User: userUUID}
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(want, nil)

	got, err := s.service.Get(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(want, got)
}

func (s *ServiceSuite) TestGetNotFound() {
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	got, err := s.service.Get(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Nil(got)
}
