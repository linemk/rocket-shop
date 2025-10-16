package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (uc *useCase) CreateOrder(ctx context.Context, info OrderInfo) (string, error) {
	if len(info.PartUUIDs) == 0 {
		return "", apperrors.ErrNoPartsSpecified
	}

	var totalPrice float32
	for _, partUUID := range info.PartUUIDs {
		partInfo, err := uc.inventoryClient.GetPart(ctx, partUUID)
		if err != nil {
			return "", apperrors.ErrPartNotFound
		}

		if partInfo.StockQuantity <= 0 {
			return "", apperrors.ErrPartOutOfStock
		}

		totalPrice += partInfo.Price
	}

	orderUUID := uuid.New()
	now := time.Now()

	order := models.Order{
		UUID:          orderUUID.String(),
		UserID:        info.UserID,
		PartUUIDs:     info.PartUUIDs,
		TotalPrice:    totalPrice,
		TransactionID: "",
		PaymentMethod: info.PaymentMethod,
		Status:        order_v1.OrderStatusPENDINGPAYMENT,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	if err := uc.orderRepository.Create(ctx, order); err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return orderUUID.String(), nil
}
