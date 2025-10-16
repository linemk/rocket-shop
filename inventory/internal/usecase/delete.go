package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
)

func (uc *useCase) DeletePart(ctx context.Context, uuid string) error {
	if err := uc.inventoryRepository.DeletePart(ctx, uuid); err != nil {
		return apperrors.ErrPartNotFound
	}
	return nil
}
