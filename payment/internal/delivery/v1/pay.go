package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func (a *API) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentUseCase.PayOrder(ctx, req.OrderUuid, req.UserUuid, req.PaymentMethod)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Payment failed: %v", err)
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
