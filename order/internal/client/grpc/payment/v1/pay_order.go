package v1

import (
	"context"
	"fmt"

	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func (c *Client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod payment_v1.PaymentMethod) (string, error) {
	resp, err := c.client.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return "", fmt.Errorf("failed to pay order: %w", err)
	}

	return resp.TransactionUuid, nil
}
