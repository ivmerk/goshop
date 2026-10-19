package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivmerk/goshop/inventory/internal/model"
)

func TestGet(t *testing.T) {
	rep := newTestRepo()
	tests := []struct {
		nameTest string
		uuid     string
		name     string
		wantErr  error
	}{
		{nameTest: "finding_details", uuid: "p1", name: "Engine"},
		{nameTest: "part_not_found", uuid: "asd", name: "Engine", wantErr: model.ErrPartNotFound},
	}
	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			part, err := rep.Get(context.Background(), tc.uuid)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Empty(t, part)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.name, part.Name)
			require.Equal(t, tc.uuid, part.UUID)
		})
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	rep := newTestRepo()
	original := rep.parts["p1"].Tags[0]
	got, err := rep.Get(context.Background(), "p1")
	require.NoError(t, err)

	got.Tags[0] = "changed"
	got.Name = "changed"

	stored := rep.parts["p1"]
	require.Equal(t, original, stored.Tags[0])
	require.NotEqual(t, "changed", stored.Name)
	for k := range got.Metadata {
		delete(got.Metadata, k)
	}
	require.NotEmpty(t, rep.parts["p1"].Metadata)

	require.Equal(t, "MB-1", *rep.parts["p1"].Metadata["series"].String)
}
