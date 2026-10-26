package part

import (
	"context"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoConverter "github.com/ivmerk/goshop/inventory/internal/repository/converter"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	cursor, err := r.collection.Find(ctx, buildFilter(filter))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var docs []*repoModel.Part
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	parts := make([]*model.Part, 0, len(docs))
	for _, d := range docs {
		parts = append(parts, repoConverter.FromRepositoryPart(d))
	}
	return parts, nil
}

func buildFilter(f model.PartsFilter) bson.M {
	filter := bson.M{}
	if len(f.UUIDs) > 0 {
		filter["_id"] = bson.M{"$in": f.UUIDs}
	}
	if len(f.Names) > 0 {
		filter["name"] = bson.M{"$in": f.Names}
	}
	if len(f.Categories) > 0 {
		filter["category"] = bson.M{"$in": f.Categories}
	}
	if len(f.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{"$in": f.ManufacturerCountries}
	}
	if len(f.Tags) > 0 {
		filter["tags"] = bson.M{"$in": f.Tags}
	}
	return filter
}
