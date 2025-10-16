package usecase

import (
	"context"
	"fmt"

	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (uc *useCase) CancelOrder(ctx context.Context, uuid string) error {
	order, err := uc.orderRepository.Get(ctx, uuid)
	if err != nil {
		return apperrors.ErrOrderNotFound
	}

	if order.Status == order_v1.OrderStatusPAID {
		return apperrors.ErrOrderAlreadyPaid
	}

	if order.Status == order_v1.OrderStatusPENDINGPAYMENT {
		updateInfo := models.OrderUpdateInfo{
			Status: &[]order_v1.OrderStatus{order_v1.OrderStatusCANCELLED}[0],
		}

		if err := uc.orderRepository.Update(ctx, uuid, updateInfo); err != nil {
			return fmt.Errorf("failed to cancel order: %w", err)
		}
	}

	return nil
}
