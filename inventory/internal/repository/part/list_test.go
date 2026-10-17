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
	}
	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			var filter model.PartsFilter
			filter.UUIDs = tc.uuid

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
			require.Equal(t, partsUUIDs, tc.partsUUIDs)
		})
	}
}
