package part

import (
	"context"

	def "github.com/ivmerk/goshop/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ def.PartRepository = (*repository)(nil)

type partsCollection interface {
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)
	InsertMany(ctx context.Context, documents any, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error)
}

type repository struct {
	collection partsCollection
}

func NewPartRepository(collection partsCollection) *repository {
	return &repository{
		collection: collection,
	}
}
