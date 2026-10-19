package part

import repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"

func newTestRepo() *repository {
	return &repository{parts: map[string]*repoModel.Part{
		"p1": {
			UUID: "p1", Name: "Engine", Category: repoModel.CategoryEngine,
			Tags:         []string{"main", "engine"},
			Manufacturer: &repoModel.Manufacturer{Country: "Germany"},
			Metadata:     map[string]*repoModel.Value{"series": {String: ptr("MB-1")}},
		},

		"p2": {
			UUID: "p2", Name: "Tank", Category: repoModel.CategoryFuel,
			Tags: []string{"fuel"}, Manufacturer: &repoModel.Manufacturer{Country: "USA"},
		},
		"p3": {
			UUID: "p3", Name: "Window", Category: repoModel.CategoryPorthole,
			Tags: []string{"glass"},
		}, // Manufacturer == nil
	}}
}
