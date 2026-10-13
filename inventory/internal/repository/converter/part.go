package converter

import (
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
		Tags:          part.Tags,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
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
		Tags:          part.Tags,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
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
