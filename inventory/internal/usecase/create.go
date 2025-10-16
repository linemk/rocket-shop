package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
)

func (uc *useCase) CreatePart(ctx context.Context, part models.Part) error {
	if err := uc.inventoryRepository.CreatePart(ctx, part); err != nil {
		return apperrors.ErrPartAlreadyExists
	}
	return nil
}
