package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
	"github.com/linemk/rocket-shop/payment/internal/repository/payment"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func TestUpdateTransaction(t *testing.T) {
	ctx := context.Background()

	repo := payment.NewRepository()

	originalTransaction := models.Transaction{
		UUID:          uuid.New().String(),
		OrderUUID:     uuid.New().String(),
		UserID:        uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
		Amount:        99.99,
		Status:        models.TransactionStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create transaction for testing
	err := repo.CreateTransaction(ctx, originalTransaction)
	require.NoError(t, err)

	// Prepare updated transaction
	updatedTransaction := originalTransaction
	updatedTransaction.Status = models.TransactionStatusCompleted
	updatedTransaction.UpdatedAt = time.Now()

	tests := []struct {
		name    string
		uuid    string
		trans   models.Transaction
		wantErr bool
	}{
		{
			name:    "successfully update transaction",
			uuid:    originalTransaction.UUID,
			trans:   updatedTransaction,
			wantErr: false,
		},
		{
			name:    "update transaction not found",
			uuid:    "nonexistent-uuid",
			trans:   updatedTransaction,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateTransaction(ctx, tt.uuid, tt.trans)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify the update
			result, err := repo.GetTransaction(ctx, tt.uuid)
			require.NoError(t, err)
			require.Equal(t, models.TransactionStatusCompleted, result.Status)
		})
	}
}
