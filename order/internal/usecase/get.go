package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

func (uc *useCase) GetOrder(ctx context.Context, uuid string) (models.Order, error) {
	return uc.orderRepository.Get(ctx, uuid)
}
