package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/linemk/rocket-shop/order/internal/usecase"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req order_v1.OptCreateOrderReq) (order_v1.CreateOrderRes, error) {
	orderInfo := usecase.OrderInfo{
		UserID:        req.Value.UserUUID.String(),
		PartUUIDs:     req.Value.PartUuids,
		PaymentMethod: order_v1.PaymentMethodPAYMENTMETHODUNSPECIFIED, // По умолчанию
	}

	orderUUID, err := a.orderUseCase.CreateOrder(ctx, orderInfo)
	if err != nil {
		return &order_v1.BadRequest{
			Code:    400,
			Message: fmt.Sprintf("Failed to create order: %v", err),
		}, nil
	}

	// Получаем созданный заказ для получения TotalPrice
	order, err := a.orderUseCase.GetOrder(ctx, orderUUID)
	if err != nil {
		return &order_v1.BadRequest{
			Code:    400,
			Message: fmt.Sprintf("Failed to get created order: %v", err),
		}, nil
	}

	orderUUIDParsed := uuid.MustParse(orderUUID)
	return &order_v1.CreateOrderResp{
		UUID:       orderUUIDParsed,
		TotalPrice: order.TotalPrice,
	}, nil
}
