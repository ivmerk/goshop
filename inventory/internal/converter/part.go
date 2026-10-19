package converter

import (
	"slices"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ivmerk/goshop/inventory/internal/model"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

func PartToProto(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return nil
	}
	out := &inventoryV1.Part{
		Uuid: part.UUID, Name: part.Name, Description: part.Description,
		Price: part.Price, StockQuantity: part.StockQuantity,
		Category:     CategoryToProto(part.Category),
		Tags:         slices.Clone(part.Tags),
		Metadata:     MetadataToProto(part.Metadata),
		CreatedAt:    timestamppb.New(part.CreatedAt),
		Dimensions:   DimensionsToProto(part.Dimensions),
		Manufacturer: ManufacturerToProto(part.Manufacturer),
	}
	if part.UpdatedAt != nil {
		out.UpdatedAt = timestamppb.New(*part.UpdatedAt)
	}
	return out
}

func PartsToProto(parts []*model.Part) []*inventoryV1.Part {
	var protoParts []*inventoryV1.Part
	for _, part := range parts {
		protoParts = append(protoParts, PartToProto(part))
	}
	return protoParts
}

func DimensionsToProto(d *model.Dimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: d.Length, Width: d.Width, Height: d.Height, Weight: d.Weight,
	}
}

func ManufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name: m.Name, Country: m.Country, Website: m.Website,
	}
}

func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	categories := make([]model.Category, 0, len(filter.GetCategories()))
	for _, c := range filter.GetCategories() {
		categories = append(categories, model.Category(c.String()))
	}

	return model.PartsFilter{
		UUIDs:                 slices.Clone(filter.GetUuids()),
		Names:                 slices.Clone(filter.GetNames()),
		Categories:            categories,
		ManufacturerCountries: slices.Clone(filter.GetManufacturerCountries()),
		Tags:                  slices.Clone(filter.GetTags()),
	}
}

func ValueToProto(v *model.Value) *inventoryV1.Value {
	var out inventoryV1.Value
	switch {
	case v.String != nil:
		out.Kind = &inventoryV1.Value_StringValue{StringValue: *v.String}
	case v.Int64 != nil:
		out.Kind = &inventoryV1.Value_Int64Value{Int64Value: *v.Int64}
	case v.Double != nil:
		out.Kind = &inventoryV1.Value_DoubleValue{DoubleValue: *v.Double}
	case v.Bool != nil:
		out.Kind = &inventoryV1.Value_BoolValue{BoolValue: *v.Bool}
	}
	return &out
}

func MetadataToProto(m map[string]*model.Value) map[string]*inventoryV1.Value {
	if m == nil {
		return nil
	}
	out := make(map[string]*inventoryV1.Value, len(m))
	for k, v := range m {
		if v == nil {
			continue
		}
		out[k] = ValueToProto(v)
	}
	return out
}

func CategoryToProto(c model.Category) inventoryV1.Category {
	return inventoryV1.Category(inventoryV1.Category_value[string(c)])
}
