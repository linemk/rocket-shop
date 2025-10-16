package v1

import (
	"github.com/linemk/rocket-shop/inventory/internal/usecase"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

type API struct {
	inventory_v1.UnimplementedInventoryServiceServer
	inventoryUseCase usecase.InventoryUseCase
}

func NewAPI(inventoryUseCase usecase.InventoryUseCase) *API {
	return &API{
		inventoryUseCase: inventoryUseCase,
	}
}
