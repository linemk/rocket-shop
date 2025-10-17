package v1

import (
	"context"

	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod payment_v1.PaymentMethod) (string, error)
	Close() error
}
