package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

// PartToProto конвертирует модель Part в protobuf Part
func PartToProto(part models.Part) *inventory_v1.Part {
	return convertToProto(part.UUID, part.Name, part.Description, part.Price, part.StockQuantity, part.Category, part.Tags, part.CreatedAt, part.UpdatedAt, part.Dimensions, part.Manufacturer)
}

// PartInfoToProto конвертирует модель PartInfo в protobuf Part
func PartInfoToProto(partInfo models.PartInfo) *inventory_v1.Part {
	return convertToProto(partInfo.UUID, partInfo.Name, partInfo.Description, partInfo.Price, partInfo.StockQuantity, partInfo.Category, partInfo.Tags, partInfo.CreatedAt, partInfo.UpdatedAt, partInfo.Dimensions, partInfo.Manufacturer)
}

// convertToProto общая функция для конвертации в protobuf Part
func convertToProto(uuid, name, description string, price float64, stockQuantity int64, category inventory_v1.Category, tags []string, createdAt, updatedAt time.Time, dimensions *models.Dimensions, manufacturer *models.Manufacturer) *inventory_v1.Part {
	protoPart := &inventory_v1.Part{
		Uuid:          uuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Tags:          tags,
		CreatedAt:     timestamppb.New(createdAt),
		UpdatedAt:     timestamppb.New(updatedAt),
	}

	// Конвертируем Dimensions
	if dimensions != nil {
		protoPart.Dimensions = &inventory_v1.Dimensions{
			Length: dimensions.Length,
			Width:  dimensions.Width,
			Height: dimensions.Height,
			Weight: dimensions.Weight,
		}
	}

	// Конвертируем Manufacturer
	if manufacturer != nil {
		protoPart.Manufacturer = &inventory_v1.Manufacturer{
			Name:    manufacturer.Name,
			Country: manufacturer.Country,
			Website: manufacturer.Website,
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
	if protoFilter == nil {
		return models.PartFilter{}
	}
	return models.PartFilter{
		UUIDs:                 protoFilter.Uuids,
		Names:                 protoFilter.Names,
		Categories:            protoFilter.Categories,
		ManufacturerCountries: protoFilter.ManufacturerCountries,
		Tags:                  protoFilter.Tags,
	}
}
