package part

import (
	"context"
	"testing"

	"github.com/ivmerk/goshop/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	rep := newTestRepo()
	tests := []struct {
		nameTest              string
		uuid                  []string
		name                  []string
		categories            []model.Category
		manufacturerCountries []string
		tags                  []string
		partsUUIDs            []string
		wantErr               error
	}{
		{nameTest: "empty_filter", partsUUIDs: []string{"p1", "p2", "p3"}},
		{nameTest: "empty_uuids", uuid: []string{"p1", "p3"}, partsUUIDs: []string{"p1", "p3"}},
		{nameTest: "filter_by_name", name: []string{"Tank"}, partsUUIDs: []string{"p2"}},
		{nameTest: "filter_by_name", categories: []model.Category{"ENGINE", "FUEL"}, partsUUIDs: []string{"p1", "p2"}},
		{nameTest: "filter_by_country", manufacturerCountries: []string{"Germany"}, partsUUIDs: []string{"p1"}},
		{nameTest: "filter_by_country_and_id", uuid: []string{"p3"}, manufacturerCountries: []string{"Germany"}},
		{nameTest: "filter_by_tags", tags: []string{"fuel", "glass"}, partsUUIDs: []string{"p2", "p3"}},
		{nameTest: "filter_by_categories_and_tags", categories: []model.Category{"ENGINE"}, tags: []string{"fuel"}},
		{nameTest: "no_matched_filter", name: []string{"NOPE"}},
	}
	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			var filter model.PartsFilter
			filter.UUIDs = tc.uuid
			filter.Names = tc.name
			filter.Categories = tc.categories
			filter.ManufacturerCountries = tc.manufacturerCountries
			filter.Tags = tc.tags

			parts, err := rep.List(context.Background(), filter)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Empty(t, parts)
				return
			}
			require.NoError(t, err)
			var partsUUIDs []string
			for _, part := range parts {
				partsUUIDs = append(partsUUIDs, part.UUID)
			}
			require.ElementsMatch(t, tc.partsUUIDs, partsUUIDs)
		})
	}
}
