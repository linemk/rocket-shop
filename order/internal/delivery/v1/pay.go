package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *order_v1.PayOrderReq, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	orderID := params.OrderUUID.String()
	paymentMethod := req.PaymentMethod

	transactionUUID, err := a.orderUseCase.PayOrder(ctx, orderID, paymentMethod)
	if err != nil {
		if err.Error() == "order not found" {
			return &order_v1.NotFoundErr{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		if err.Error() == "order cannot be paid in current status" {
			return &order_v1.ConflictErr{
				Code:    409,
				Message: "Order cannot be paid in current status",
			}, nil
		}
		return &order_v1.BadRequest{
			Code:    400,
			Message: fmt.Sprintf("Payment failed: %v", err),
		}, nil
	}

	transactionUUIDParsed := uuid.MustParse(transactionUUID)
	return &order_v1.PayOrderResp{
		TransactionUUID: transactionUUIDParsed,
	}, nil
}
