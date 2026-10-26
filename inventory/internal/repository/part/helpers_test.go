package part

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

const (
	partUUID1 = "11111111-1111-1111-1111-111111111111"
	partUUID2 = "22222222-2222-2222-2222-222222222222"
)

// fakeCollection подменяет *mongo.Collection: каждый тест задаёт нужные ему функции.
// Вызов метода без заданной функции — ошибка теста (неожиданное обращение к БД).
type fakeCollection struct {
	t *testing.T

	findOne    func(filter any) *mongo.SingleResult
	find       func(filter any) (*mongo.Cursor, error)
	count      func(filter any) (int64, error)
	insertMany func(documents any) (*mongo.InsertManyResult, error)
}

var _ partsCollection = (*fakeCollection)(nil)

func (f *fakeCollection) FindOne(_ context.Context, filter any, _ ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	if f.findOne == nil {
		f.t.Fatal("неожиданный вызов FindOne")
	}
	return f.findOne(filter)
}

func (f *fakeCollection) Find(_ context.Context, filter any, _ ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	if f.find == nil {
		f.t.Fatal("неожиданный вызов Find")
	}
	return f.find(filter)
}

func (f *fakeCollection) CountDocuments(_ context.Context, filter any, _ ...options.Lister[options.CountOptions]) (int64, error) {
	if f.count == nil {
		f.t.Fatal("неожиданный вызов CountDocuments")
	}
	return f.count(filter)
}

func (f *fakeCollection) InsertMany(_ context.Context, documents any, _ ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error) {
	if f.insertMany == nil {
		f.t.Fatal("неожиданный вызов InsertMany")
	}
	return f.insertMany(documents)
}

func newTestRepo(t *testing.T) (*repository, *fakeCollection) {
	t.Helper()

	coll := &fakeCollection{t: t}
	return NewPartRepository(coll), coll
}

func testPart() *repoModel.Part {
	return &repoModel.Part{
		UUID:         partUUID1,
		Name:         "Engine",
		Category:     repoModel.CategoryEngine,
		Tags:         []string{"main", "engine"},
		Manufacturer: &repoModel.Manufacturer{Country: "Germany"},
		Metadata:     map[string]*repoModel.Value{"series": {String: ptr("MB-1")}},
	}
}
