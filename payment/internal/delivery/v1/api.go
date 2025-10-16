package v1

import (
	"github.com/linemk/rocket-shop/payment/internal/usecase"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

type API struct {
	payment_v1.UnimplementedPaymentServiceServer
	paymentUseCase usecase.PaymentUseCase
}

func NewAPI(paymentUseCase usecase.PaymentUseCase) *API {
	return &API{
		paymentUseCase: paymentUseCase,
	}
}
