package app

import (
	"context"

	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"

	v1 "github.com/linemk/rocket-shop/payment/internal/delivery/v1"
	"github.com/linemk/rocket-shop/payment/internal/repository"
	paymentRepository "github.com/linemk/rocket-shop/payment/internal/repository/payment"
	"github.com/linemk/rocket-shop/payment/internal/usecase"
)

type diContainer struct {
	paymentV1API payment_v1.PaymentServiceServer

	paymentUseCase usecase.PaymentUseCase

	paymentRepository repository.PaymentRepository
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) payment_v1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = v1.NewAPI(d.PaymentUseCase(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentUseCase(ctx context.Context) usecase.PaymentUseCase {
	if d.paymentUseCase == nil {
		d.paymentUseCase = usecase.NewUseCase(d.PaymentRepository(ctx))
	}

	return d.paymentUseCase
}

func (d *diContainer) PaymentRepository(ctx context.Context) repository.PaymentRepository {
	if d.paymentRepository == nil {
		d.paymentRepository = paymentRepository.NewRepository()
	}

	return d.paymentRepository
}
