package models

import (
	"github.com/google/uuid"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
	"time"
)

type OrderUpdateInfo struct {
	Status        *order_v1.OrderStatus
	TransactionID *string
	PaymentMethod *order_v1.PaymentMethod
}

type Order struct {
	UUID          string                 `json:"order_id"`
	UserID        string                 `json:"user_id"`
	PartUUIDs     []uuid.UUID            `json:"details_id"`
	TotalPrice    float32                `json:"total_price"`
	TransactionID string                 `json:"transaction_id"`
	PaymentMethod order_v1.PaymentMethod `json:"payment_method"`
	Status        order_v1.OrderStatus   `json:"status"`
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
