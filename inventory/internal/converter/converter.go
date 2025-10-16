package converter

import (
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PartToProto конвертирует модель Part в protobuf Part
func PartToProto(part models.Part) *inventory_v1.Part {
	protoPart := &inventory_v1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Tags:          part.Tags,
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}

	// Конвертируем Dimensions
	if part.Dimensions != nil {
		protoPart.Dimensions = &inventory_v1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	// Конвертируем Manufacturer
	if part.Manufacturer != nil {
		protoPart.Manufacturer = &inventory_v1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	return protoPart
}

// PartInfoToProto конвертирует модель PartInfo в protobuf Part
func PartInfoToProto(partInfo models.PartInfo) *inventory_v1.Part {
	protoPart := &inventory_v1.Part{
		Uuid:          partInfo.UUID,
		Name:          partInfo.Name,
		Description:   partInfo.Description,
		Price:         partInfo.Price,
		StockQuantity: partInfo.StockQuantity,
		Category:      partInfo.Category,
		Tags:          partInfo.Tags,
		CreatedAt:     timestamppb.New(partInfo.CreatedAt),
		UpdatedAt:     timestamppb.New(partInfo.UpdatedAt),
	}

	// Конвертируем Dimensions
	if partInfo.Dimensions != nil {
		protoPart.Dimensions = &inventory_v1.Dimensions{
			Length: partInfo.Dimensions.Length,
			Width:  partInfo.Dimensions.Width,
			Height: partInfo.Dimensions.Height,
			Weight: partInfo.Dimensions.Weight,
		}
	}

	// Конвертируем Manufacturer
	if partInfo.Manufacturer != nil {
		protoPart.Manufacturer = &inventory_v1.Manufacturer{
			Name:    partInfo.Manufacturer.Name,
			Country: partInfo.Manufacturer.Country,
			Website: partInfo.Manufacturer.Website,
		}
	}

	return protoPart
}

// ProtoToPart конвертирует protobuf Part в модель Part
func ProtoToPart(protoPart *inventory_v1.Part) models.Part {
	part := models.Part{
		UUID:          protoPart.Uuid,
		Name:          protoPart.Name,
		Description:   protoPart.Description,
		Price:         protoPart.Price,
		StockQuantity: protoPart.StockQuantity,
		Category:      protoPart.Category,
		Tags:          protoPart.Tags,
		CreatedAt:     protoPart.CreatedAt.AsTime(),
		UpdatedAt:     protoPart.UpdatedAt.AsTime(),
	}

	// Конвертируем Dimensions
	if protoPart.Dimensions != nil {
		part.Dimensions = &models.Dimensions{
			Length: protoPart.Dimensions.Length,
			Width:  protoPart.Dimensions.Width,
			Height: protoPart.Dimensions.Height,
			Weight: protoPart.Dimensions.Weight,
		}
	}

	// Конвертируем Manufacturer
	if protoPart.Manufacturer != nil {
		part.Manufacturer = &models.Manufacturer{
			Name:    protoPart.Manufacturer.Name,
			Country: protoPart.Manufacturer.Country,
			Website: protoPart.Manufacturer.Website,
		}
	}

	return part
}

// ProtoToPartFilter конвертирует protobuf PartsFilter в модель PartFilter
func ProtoToPartFilter(protoFilter *inventory_v1.PartsFilter) models.PartFilter {
	return models.PartFilter{
		UUIDs:                 protoFilter.Uuids,
		Names:                 protoFilter.Names,
		Categories:            protoFilter.Categories,
		ManufacturerCountries: protoFilter.ManufacturerCountries,
		Tags:                  protoFilter.Tags,
	}
}
