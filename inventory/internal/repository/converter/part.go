package converter

import (
	"slices"

	"github.com/ivmerk/goshop/inventory/internal/model"
	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

func ToRepositoryPart(part *model.Part) *repoModel.Part {
	if part == nil {
		return nil
	}
	return &repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		StockQuantity: part.StockQuantity,
		Price:         part.Price,
		Category:      repoModel.Category(part.Category),
		Tags:          slices.Clone(part.Tags),
		Metadata:      ToRepositoryMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     clonePtr(part.UpdatedAt),
		Dimensions:    ToRepositoryDimensions(part.Dimensions),
		Manufacturer:  ToRepositoryManufacturer(part.Manufacturer),
	}
}

func ToRepositoryDimensions(dimensions *model.Dimensions) *repoModel.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToRepositoryManufacturer(manufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToModelPart(part *repoModel.Part) *model.Part {
	if part == nil {
		return nil
	}
	return &model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		StockQuantity: part.StockQuantity,
		Price:         part.Price,
		Category:      model.Category(part.Category),
		Tags:          slices.Clone(part.Tags),
		Metadata:      ToModelMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     clonePtr(part.UpdatedAt),
		Dimensions:    ToModelDimensions(part.Dimensions),
		Manufacturer:  ToModelManufacturer(part.Manufacturer),
	}
}

func ToModelDimensions(dimensions *repoModel.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToModelManufacturer(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func clonePtr[T any](s *T) *T {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}

func ToRepositoryValue(v *model.Value) *repoModel.Value {
	if v == nil {
		return nil
	}
	return &repoModel.Value{
		String: clonePtr(v.String),
		Int64:  clonePtr(v.Int64),
		Double: clonePtr(v.Double),
		Bool:   clonePtr(v.Bool),
	}
}

func ToModelValue(v *repoModel.Value) *model.Value {
	if v == nil {
		return nil
	}
	return &model.Value{
		String: clonePtr(v.String),
		Int64:  clonePtr(v.Int64),
		Double: clonePtr(v.Double),
		Bool:   clonePtr(v.Bool),
	}
}

func ToRepositoryMetadata(m map[string]*model.Value) map[string]*repoModel.Value {
	if m == nil {
		return nil
	}
	out := make(map[string]*repoModel.Value, len(m))
	for k, v := range m {
		if v == nil {
			continue
		}
		out[k] = ToRepositoryValue(v)
	}
	return out
}

func ToModelMetadata(m map[string]*repoModel.Value) map[string]*model.Value {
	if m == nil {
		return nil
	}
	out := make(map[string]*model.Value, len(m))
	for k, v := range m {
		if v == nil {
			continue
		}
		out[k] = ToModelValue(v)
	}
	return out
}
