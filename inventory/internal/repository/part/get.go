package part

import (
	"context"
	"errors"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoConverter "github.com/ivmerk/goshop/inventory/internal/repository/converter"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {

	var p repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}

	return repoConverter.FromRepositoryPart(&p), nil
}
