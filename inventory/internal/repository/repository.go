package repository

import (
	"context"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
)

// InventoryRepository определяет интерфейс для работы с инвентарем
type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (models.Part, error)
	ListParts(ctx context.Context, filter models.PartFilter) ([]models.Part, error)
	CreatePart(ctx context.Context, part models.Part) error
	UpdatePart(ctx context.Context, uuid string, part models.Part) error
	DeletePart(ctx context.Context, uuid string) error
}
