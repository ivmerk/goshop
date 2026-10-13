package part

import (
	"github.com/ivmerk/goshop/inventory/internal/repository"
	def "github.com/ivmerk/goshop/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	repository repository.PartRepository
}

func NewPartService(partRepository repository.PartRepository) *service {
	return &service{
		repository: partRepository}
}
