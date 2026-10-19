package part

import (
	"context"
	"testing"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tests := []struct {
		name   string
		filter model.PartsFilter
		want   []string
	}{
		{
			name: "empty_filter_returns_all",
			want: []string{"p1", "p2", "p3"},
		},
		{
			name:   "filter_by_uuids",
			filter: model.PartsFilter{UUIDs: []string{"p1", "p3"}},
			want:   []string{"p1", "p3"},
		},
		{
			name:   "filter_by_names",
			filter: model.PartsFilter{Names: []string{"Tank"}},
			want:   []string{"p2"},
		},
		{
			name:   "filter_by_categories",
			filter: model.PartsFilter{Categories: []model.Category{model.CategoryEngine, model.CategoryFuel}},
			want:   []string{"p1", "p2"},
		},
		{
			name:   "filter_by_country",
			filter: model.PartsFilter{ManufacturerCountries: []string{"Germany"}},
			want:   []string{"p1"},
		},
		{
			name:   "country_filter_skips_part_without_manufacturer",
			filter: model.PartsFilter{UUIDs: []string{"p3"}, ManufacturerCountries: []string{"Germany"}},
			want:   []string{},
		},
		{
			name:   "filter_by_tags_matches_any_tag",
			filter: model.PartsFilter{Tags: []string{"fuel", "glass"}},
			want:   []string{"p2", "p3"},
		},
		{
			name: "fields_are_combined_with_and",
			filter: model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine},
				Tags:       []string{"engine"},
			},
			want: []string{"p1"},
		},
		{
			name: "conflicting_fields_return_nothing",
			filter: model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine},
				Tags:       []string{"fuel"},
			},
			want: []string{},
		},
		{
			name:   "no_match_returns_empty_list",
			filter: model.PartsFilter{Names: []string{"NOPE"}},
			want:   []string{},
		},
	}

	rep := newTestRepo()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts, err := rep.List(context.Background(), tc.filter)
			require.NoError(t, err)

			got := make([]string, 0, len(parts))
			for _, part := range parts {
				got = append(got, part.UUID)
			}
			require.ElementsMatch(t, tc.want, got)
		})
	}
}

func TestList_EmptyRepository(t *testing.T) {
	rep := &repository{parts: map[string]*repoModel.Part{}}

	parts, err := rep.List(context.Background(), model.PartsFilter{})

	require.NoError(t, err)
	require.Empty(t, parts)
}

func TestNewPartRepository_IsSeeded(t *testing.T) {
	rep := NewPartRepository()

	parts, err := rep.List(context.Background(), model.PartsFilter{})
	require.NoError(t, err)
	require.Len(t, parts, 4)

	withMetadata := 0
	for _, part := range parts {
		require.NotEmpty(t, part.UUID)
		require.NotEmpty(t, part.Name)
		if len(part.Metadata) > 0 {
			withMetadata++
		}
	}
	require.Equal(t, 1, withMetadata)
}
