package part

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

func TestBuildFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter model.PartsFilter
		want   bson.M
	}{
		{
			name: "empty_filter",
			want: bson.M{},
		},
		{
			name:   "uuids",
			filter: model.PartsFilter{UUIDs: []string{partUUID1, partUUID2}},
			want:   bson.M{"_id": bson.M{"$in": []string{partUUID1, partUUID2}}},
		},
		{
			name:   "names",
			filter: model.PartsFilter{Names: []string{"Engine"}},
			want:   bson.M{"name": bson.M{"$in": []string{"Engine"}}},
		},
		{
			name:   "categories",
			filter: model.PartsFilter{Categories: []model.Category{model.CategoryEngine, model.CategoryFuel}},
			want:   bson.M{"category": bson.M{"$in": []model.Category{model.CategoryEngine, model.CategoryFuel}}},
		},
		{
			name:   "manufacturer_countries",
			filter: model.PartsFilter{ManufacturerCountries: []string{"Germany"}},
			want:   bson.M{"manufacturer.country": bson.M{"$in": []string{"Germany"}}},
		},
		{
			name:   "tags",
			filter: model.PartsFilter{Tags: []string{"main", "glass"}},
			want:   bson.M{"tags": bson.M{"$in": []string{"main", "glass"}}},
		},
		{
			name: "all_fields_are_combined",
			filter: model.PartsFilter{
				UUIDs:                 []string{partUUID1},
				Names:                 []string{"Engine"},
				Categories:            []model.Category{model.CategoryEngine},
				ManufacturerCountries: []string{"Germany"},
				Tags:                  []string{"main"},
			},
			want: bson.M{
				"_id":                  bson.M{"$in": []string{partUUID1}},
				"name":                 bson.M{"$in": []string{"Engine"}},
				"category":             bson.M{"$in": []model.Category{model.CategoryEngine}},
				"manufacturer.country": bson.M{"$in": []string{"Germany"}},
				"tags":                 bson.M{"$in": []string{"main"}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, buildFilter(tc.filter))
		})
	}
}

func TestList_Success(t *testing.T) {
	rep, coll := newTestRepo(t)
	filter := model.PartsFilter{Categories: []model.Category{model.CategoryEngine}}
	second := &repoModel.Part{UUID: partUUID2, Name: "Tank", Category: repoModel.CategoryFuel}

	coll.find = func(got any) (*mongo.Cursor, error) {
		require.Equal(t, buildFilter(filter), got)
		return mongo.NewCursorFromDocuments([]any{testPart(), second}, nil, nil)
	}

	parts, err := rep.List(context.Background(), filter)

	require.NoError(t, err)
	require.Len(t, parts, 2)
	require.Equal(t, partUUID1, parts[0].UUID)
	require.Equal(t, model.CategoryEngine, parts[0].Category)
	require.Equal(t, "Germany", parts[0].Manufacturer.Country)
	require.Equal(t, partUUID2, parts[1].UUID)
}

func TestList_Empty(t *testing.T) {
	rep, coll := newTestRepo(t)
	coll.find = func(any) (*mongo.Cursor, error) {
		return mongo.NewCursorFromDocuments([]any{}, nil, nil)
	}

	parts, err := rep.List(context.Background(), model.PartsFilter{})

	require.NoError(t, err)
	require.NotNil(t, parts)
	require.Empty(t, parts)
}

func TestList_FindError(t *testing.T) {
	rep, coll := newTestRepo(t)
	dbErr := errors.New("connection lost")
	coll.find = func(any) (*mongo.Cursor, error) {
		return nil, dbErr
	}

	parts, err := rep.List(context.Background(), model.PartsFilter{})

	require.ErrorIs(t, err, dbErr)
	require.Nil(t, parts)
}

func TestList_DecodeError(t *testing.T) {
	rep, coll := newTestRepo(t)
	coll.find = func(any) (*mongo.Cursor, error) {
		// _id числом не декодируется в строковый UUID — cursor.All вернёт ошибку
		return mongo.NewCursorFromDocuments([]any{bson.M{"_id": 123}}, nil, nil)
	}

	parts, err := rep.List(context.Background(), model.PartsFilter{})

	require.Error(t, err)
	require.Nil(t, parts)
}
