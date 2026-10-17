package part

import "github.com/ivmerk/goshop/inventory/internal/model"

func (s *ServiceSuite) TestGetSuccess() {
	want := &model.Part{UUID: "p1", Name: "Engine"}
	s.partRepository.On("Get", s.ctx, "p1").Return(want, nil)

	got, err := s.service.Get(s.ctx, "p1")

	s.Require().NoError(err)
	s.Require().Equal(want, got)
}

func (s *ServiceSuite) TestGetNotFound() {
	s.partRepository.On("Get", s.ctx, "p1").Return(nil, model.ErrPartNotFound)

	got, err := s.service.Get(s.ctx, "p1")

	s.Require().ErrorIs(err, model.ErrPartNotFound)
	s.Require().Nil(got)
}
