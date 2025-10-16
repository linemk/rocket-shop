package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
)

func (uc *useCase) ListTransactions(ctx context.Context, orderUUID string) ([]models.Transaction, error) {
	return uc.paymentRepository.ListTransactions(ctx, orderUUID)
}
