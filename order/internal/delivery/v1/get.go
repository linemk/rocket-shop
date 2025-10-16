package v1

import (
	"context"

	"github.com/google/uuid"

	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	orderUUID := params.OrderUUID.String()

	order, err := a.orderUseCase.GetOrder(ctx, orderUUID)
	if err != nil {
		return &order_v1.NotFoundErr{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	var transactionUUID uuid.UUID
	if order.TransactionID != "" {
		transactionUUID = uuid.MustParse(order.TransactionID)
	}

	response := &order_v1.GetOrderResp{
		OrderUUID:       params.OrderUUID,
		UserUUID:        uuid.MustParse(order.UserID),
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}

	return response, nil
}
