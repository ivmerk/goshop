package part

import (
	"context"
	"testing"

	"github.com/ivmerk/goshop/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx            context.Context
	partRepository *mocks.PartRepository
	service        *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.partRepository = mocks.NewPartRepository(s.T())
	s.service = NewPartService(s.partRepository)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
