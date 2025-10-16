package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
	"github.com/linemk/rocket-shop/payment/internal/repository"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

type PaymentUseCase interface {
	PayOrder(ctx context.Context, orderUUID, userID string, paymentMethod payment_v1.PaymentMethod) (string, error)
	GetTransaction(ctx context.Context, transactionUUID string) (models.Transaction, error)
	ListTransactions(ctx context.Context, orderUUID string) ([]models.Transaction, error)
}

type useCase struct {
	paymentRepository repository.PaymentRepository
}

func NewUseCase(paymentRepository repository.PaymentRepository) PaymentUseCase {
	return &useCase{
		paymentRepository: paymentRepository,
	}
}
