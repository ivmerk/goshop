package part

import (
	"context"
	"slices"

	"github.com/ivmerk/goshop/inventory/internal/model"

	repoConverter "github.com/ivmerk/goshop/inventory/internal/repository/converter"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []*model.Part
	for _, part := range r.parts {
		if !matches(part, filter) {
			continue
		}
		parts = append(parts, repoConverter.FromRepositoryPart(part))
	}

	return parts, nil
}

func matches(p *repoModel.Part, f model.PartsFilter) bool {
	if len(f.UUIDs) > 0 && !slices.Contains(f.UUIDs, p.UUID) {
		return false
	}
	if len(f.Names) > 0 && !slices.Contains(f.Names, p.Name) {
		return false
	}
	if len(f.Categories) > 0 && !slices.Contains(f.Categories, model.Category(p.Category)) {
		return false
	}
	if len(f.ManufacturerCountries) > 0 &&
		(p.Manufacturer == nil || !slices.Contains(f.ManufacturerCountries, p.Manufacturer.Country)) {
		return false
	}
	if len(f.Tags) > 0 && !slices.ContainsFunc(p.Tags, func(t string) bool { return slices.Contains(f.Tags, t) }) {
		return false
	}
	return true
}
