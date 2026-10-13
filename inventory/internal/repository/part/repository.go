package part

import (
	"sync"

	def "github.com/ivmerk/goshop/inventory/internal/repository"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

func NewPartRepository() *repository {
	r := &repository{
		parts: make(map[string]*repoModel.Part),
	}
	r.seed()
	return r
}
