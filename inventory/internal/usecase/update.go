package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
)

func (uc *useCase) UpdatePart(ctx context.Context, uuid string, part models.Part) error {
	if err := uc.inventoryRepository.UpdatePart(ctx, uuid, part); err != nil {
		return apperrors.ErrPartNotFound
	}
	return nil
}
