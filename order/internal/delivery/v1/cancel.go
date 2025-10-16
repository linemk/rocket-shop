package v1

import (
	"context"
	"fmt"

	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	orderID := params.OrderUUID.String()

	err := a.orderUseCase.CancelOrder(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return &order_v1.NotFoundErr{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		if err.Error() == "order already paid and cannot be cancelled" {
			return &order_v1.ConflictErr{
				Code:    409,
				Message: "Order already paid and cannot be cancelled",
			}, nil
		}
		// Для других ошибок возвращаем ConflictErr
		return &order_v1.ConflictErr{
			Code:    409,
			Message: fmt.Sprintf("Failed to cancel order: %v", err),
		}, nil
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
