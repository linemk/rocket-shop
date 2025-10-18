package mocks

//go:generate mockgen --package mocks --destination payment_repository_mock.go github.com/linemk/rocket-shop/payment/internal/repository PaymentRepository
//go:generate mockgen --package mocks --destination payment_usecase_mock.go github.com/linemk/rocket-shop/payment/internal/usecase PaymentUseCase
