package converter

import (
	"testing"
	"time"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
	"github.com/stretchr/testify/require"
)

func domainPart() *model.Part {
	updated := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	s, i, d, b := "MB-1", int64(42), 1200.5, true
	return &model.Part{
		UUID: "p1", Name: "Engine", Description: "desc", Price: 12500.5, StockQuantity: 5,
		Category: model.CategoryEngine, Tags: []string{"main", "engine"},
		Metadata: map[string]*model.Value{
			"series": {String: &s}, "batch": {Int64: &i}, "thrust": {Double: &d}, "tested": {Bool: &b},
		},
		CreatedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: &updated,
		Dimensions:   &model.Dimensions{Length: 300, Width: 120, Height: 120, Weight: 850},
		Manufacturer: &model.Manufacturer{Name: "Rocket GmbH", Country: "Germany", Website: "https://rocket.example"},
	}
}

func TestPartRoundTrip(t *testing.T) {
	want := domainPart()

	got := FromRepositoryPart(ToRepositoryPart(want))

	require.Equal(t, want, got)
}

func TestNilInputs(t *testing.T) {
	require.Nil(t, ToRepositoryPart(nil))
	require.Nil(t, FromRepositoryPart(nil))
	require.Nil(t, ToRepositoryDimensions(nil))
	require.Nil(t, ToModelDimensions(nil))
	require.Nil(t, ToRepositoryManufacturer(nil))
	require.Nil(t, ToModelManufacturer(nil))
	require.Nil(t, ToRepositoryValue(nil))
	require.Nil(t, ToModelValue(nil))
	require.Nil(t, ToRepositoryMetadata(nil))
	require.Nil(t, ToModelMetadata(nil))
}

func TestMetadataSkipsNilValues(t *testing.T) {
	s := "x"

	toRepo := ToRepositoryMetadata(map[string]*model.Value{"a": {String: &s}, "b": nil})
	toModel := ToModelMetadata(map[string]*repoModel.Value{"a": {String: &s}, "b": nil})

	require.Len(t, toRepo, 1)
	require.Contains(t, toRepo, "a")
	require.Len(t, toModel, 1)
	require.Contains(t, toModel, "a")
}

func TestToRepositoryPartCopiesData(t *testing.T) {
	p := domainPart()

	got := ToRepositoryPart(p)
	p.Tags[0] = "changed"
	*p.Metadata["series"].String = "changed"
	*p.UpdatedAt = time.Time{}

	require.Equal(t, "main", got.Tags[0])
	require.Equal(t, "MB-1", *got.Metadata["series"].String)
	require.False(t, got.UpdatedAt.IsZero())
}

func TestFromRepositoryPartCopiesData(t *testing.T) {
	repoPart := ToRepositoryPart(domainPart())

	got := FromRepositoryPart(repoPart)
	repoPart.Tags[0] = "changed"
	*repoPart.Metadata["series"].String = "changed"
	*repoPart.UpdatedAt = time.Time{}

	require.Equal(t, "main", got.Tags[0])
	require.Equal(t, "MB-1", *got.Metadata["series"].String)
	require.False(t, got.UpdatedAt.IsZero())
}
