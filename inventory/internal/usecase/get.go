package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
)

func (uc *useCase) GetPart(ctx context.Context, uuid string) (models.PartInfo, error) {
	part, err := uc.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return models.PartInfo{}, apperrors.ErrPartNotFound
	}

	// Конвертируем Part в PartInfo
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

	return partInfo, nil
}
