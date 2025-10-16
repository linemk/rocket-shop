package v1

import (
	"github.com/linemk/rocket-shop/order/internal/usecase"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

type api struct {
	order_v1.UnimplementedHandler

	orderUseCase usecase.OrderUseCase
}

func NewAPI(orderUseCase usecase.OrderUseCase) *api {
	return &api{
		orderUseCase: orderUseCase,
	}
}
