package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
)

func (uc *useCase) ListParts(ctx context.Context, filter models.PartFilter) ([]models.PartInfo, error) {
	parts, err := uc.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Конвертируем []Part в []PartInfo
	partInfos := make([]models.PartInfo, 0, len(parts))
	for _, part := range parts {
		partInfo := models.PartInfo{
			UUID:          part.UUID,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      part.Category,
			Dimensions:    part.Dimensions,
			Manufacturer:  part.Manufacturer,
			Tags:          part.Tags,
			Metadata:      part.Metadata,
			CreatedAt:     part.CreatedAt,
			UpdatedAt:     part.UpdatedAt,
		}
		partInfos = append(partInfos, partInfo)
	}

	return partInfos, nil
}
