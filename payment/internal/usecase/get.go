package usecase

import (
	"context"

	"github.com/linemk/rocket-shop/payment/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
)

func (uc *useCase) GetTransaction(ctx context.Context, transactionUUID string) (models.Transaction, error) {
	transaction, err := uc.paymentRepository.GetTransaction(ctx, transactionUUID)
	if err != nil {
		return models.Transaction{}, apperrors.ErrTransactionNotFound
	}
	return transaction, nil
}
