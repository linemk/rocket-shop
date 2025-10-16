package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	"github.com/linemk/rocket-shop/inventory/internal/repository"
)

type InventoryUseCase interface {
	GetPart(ctx context.Context, uuid string) (models.PartInfo, error)
	ListParts(ctx context.Context, filter models.PartFilter) ([]models.PartInfo, error)
	CreatePart(ctx context.Context, part models.Part) error
	UpdatePart(ctx context.Context, uuid string, part models.Part) error
	DeletePart(ctx context.Context, uuid string) error
}

type useCase struct {
	inventoryRepository repository.InventoryRepository
}

func NewUseCase(inventoryRepository repository.InventoryRepository) InventoryUseCase {
	return &useCase{
		inventoryRepository: inventoryRepository,
	}
}
