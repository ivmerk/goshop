package part

import (
	"errors"

	"github.com/ivmerk/goshop/inventory/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
	filter := model.PartsFilter{
		Categories: []model.Category{model.CategoryEngine},
		Tags:       []string{"main"},
	}
	want := []*model.Part{
		{UUID: "p1", Name: "Engine"},
		{UUID: "p2", Name: "Booster"},
	}
	s.partRepository.On("List", s.ctx, filter).Return(want, nil)

	got, err := s.service.List(s.ctx, filter)

	s.Require().NoError(err)
	s.Require().Equal(want, got)
}

func (s *ServiceSuite) TestListEmptyResult() {
	filter := model.PartsFilter{Names: []string{"NOPE"}}
	s.partRepository.On("List", s.ctx, filter).Return([]*model.Part{}, nil)

	got, err := s.service.List(s.ctx, filter)

	s.Require().NoError(err)
	s.Require().Empty(got)
}

func (s *ServiceSuite) TestListRepositoryError() {
	filter := model.PartsFilter{}
	repoErr := errors.New("storage failure")
	s.partRepository.On("List", s.ctx, filter).Return(nil, repoErr)

	got, err := s.service.List(s.ctx, filter)

	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(got)
}
