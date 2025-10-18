package mocks

//go:generate mockgen --package mocks --destination inventory_repository_mock.go github.com/linemk/rocket-shop/inventory/internal/repository InventoryRepository
//go:generate mockgen --package mocks --destination inventory_usecase_mock.go github.com/linemk/rocket-shop/inventory/internal/usecase InventoryUseCase
