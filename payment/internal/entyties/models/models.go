package models

import (
	"time"

	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

// Transaction представляет транзакцию платежа
type Transaction struct {
	UUID          string
	OrderUUID     string
	UserID        string
	PaymentMethod payment_v1.PaymentMethod
	Amount        float64
	Status        TransactionStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TransactionStatus представляет статус транзакции
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

// PaymentRequest представляет запрос на платеж
type PaymentRequest struct {
	OrderUUID     string
	UserID        string
	PaymentMethod payment_v1.PaymentMethod
	Amount        float64
}

// PaymentResponse представляет ответ на платеж
type PaymentResponse struct {
	TransactionUUID string
	Status          TransactionStatus
}
