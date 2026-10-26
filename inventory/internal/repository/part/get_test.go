package part

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/ivmerk/goshop/inventory/internal/model"
)

func TestGet_Success(t *testing.T) {
	rep, coll := newTestRepo(t)
	coll.findOne = func(filter any) *mongo.SingleResult {
		require.Equal(t, bson.M{"_id": partUUID1}, filter)
		return mongo.NewSingleResultFromDocument(testPart(), nil, nil)
	}

	got, err := rep.Get(context.Background(), partUUID1)

	require.NoError(t, err)
	require.Equal(t, partUUID1, got.UUID)
	require.Equal(t, "Engine", got.Name)
	require.Equal(t, model.CategoryEngine, got.Category)
}

func TestGet_NotFound(t *testing.T) {
	rep, coll := newTestRepo(t)
	coll.findOne = func(any) *mongo.SingleResult {
		return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	}

	got, err := rep.Get(context.Background(), partUUID1)

	require.ErrorIs(t, err, model.ErrPartNotFound)
	require.Nil(t, got)
}

func TestGet_DBError(t *testing.T) {
	rep, coll := newTestRepo(t)
	dbErr := errors.New("connection lost")
	coll.findOne = func(any) *mongo.SingleResult {
		return mongo.NewSingleResultFromDocument(bson.D{}, dbErr, nil)
	}

	got, err := rep.Get(context.Background(), partUUID1)

	require.ErrorIs(t, err, dbErr)
	require.NotErrorIs(t, err, model.ErrPartNotFound)
	require.Nil(t, got)
}
