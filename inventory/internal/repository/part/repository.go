package part

import (
	def "github.com/ivmerk/goshop/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewPartRepository(collection *mongo.Collection) *repository {
	r := &repository{
		collection: collection,
	}
	return r
}
