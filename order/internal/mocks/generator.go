package mocks

//go:generate mockgen --package mocks --destination order_repository_mock.go github.com/linemk/rocket-shop/order/internal/repository OrderRepository
//go:generate mockgen --package mocks --destination inventory_client_mock.go github.com/linemk/rocket-shop/order/internal/client/grpc/inventory/v1 InventoryClient
//go:generate mockgen --package mocks --destination payment_client_mock.go github.com/linemk/rocket-shop/order/internal/client/grpc/payment/v1 PaymentClient
//go:generate mockgen --package mocks --destination order_usecase_mock.go github.com/linemk/rocket-shop/order/internal/usecase OrderUseCase
