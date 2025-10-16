package repository

import (
	"context"

	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
)

// PaymentRepository определяет интерфейс для работы с платежами
type PaymentRepository interface {
	CreateTransaction(ctx context.Context, transaction models.Transaction) error
	GetTransaction(ctx context.Context, uuid string) (models.Transaction, error)
	UpdateTransaction(ctx context.Context, uuid string, transaction models.Transaction) error
	ListTransactions(ctx context.Context, orderUUID string) ([]models.Transaction, error)
}
