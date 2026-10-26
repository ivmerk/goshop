package part

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

func TestSeed_AlreadyFilled(t *testing.T) {
	rep, coll := newTestRepo(t)
	coll.count = func(filter any) (int64, error) {
		require.Equal(t, bson.M{}, filter)
		return 4, nil
	}
	// insertMany не задан: вызов вставки уронит тест через fakeCollection

	require.NoError(t, rep.Seed(context.Background()))
}

func TestSeed_EmptyCollection(t *testing.T) {
	rep, coll := newTestRepo(t)
	coll.count = func(any) (int64, error) { return 0, nil }
	inserted := false
	coll.insertMany = func(documents any) (*mongo.InsertManyResult, error) {
		inserted = true

		parts, ok := documents.([]*repoModel.Part)
		require.True(t, ok)
		require.Len(t, parts, 4)

		uuids := make(map[string]struct{}, len(parts))
		for _, p := range parts {
			require.NotEmpty(t, p.UUID)
			require.NotEmpty(t, p.Name)
			require.False(t, p.CreatedAt.IsZero())
			uuids[p.UUID] = struct{}{}
		}
		require.Len(t, uuids, 4, "UUID деталей должны быть уникальными")

		return &mongo.InsertManyResult{}, nil
	}

	require.NoError(t, rep.Seed(context.Background()))
	require.True(t, inserted)
}

func TestSeed_CountError(t *testing.T) {
	rep, coll := newTestRepo(t)
	dbErr := errors.New("connection lost")
	coll.count = func(any) (int64, error) { return 0, dbErr }

	err := rep.Seed(context.Background())

	require.ErrorIs(t, err, dbErr)
}

func TestSeed_InsertError(t *testing.T) {
	rep, coll := newTestRepo(t)
	dbErr := errors.New("write failed")
	coll.count = func(any) (int64, error) { return 0, nil }
	coll.insertMany = func(any) (*mongo.InsertManyResult, error) { return nil, dbErr }

	err := rep.Seed(context.Background())

	require.ErrorIs(t, err, dbErr)
}
