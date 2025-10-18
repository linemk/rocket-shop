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

func TestGetTransaction(t *testing.T) {
	ctx := context.Background()

	repo := payment.NewRepository()

	transaction := models.Transaction{
		UUID:          uuid.New().String(),
		OrderUUID:     uuid.New().String(),
		UserID:        uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
		Amount:        99.99,
		Status:        models.TransactionStatusCompleted,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create transaction for testing
	err := repo.CreateTransaction(ctx, transaction)
	require.NoError(t, err)

	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{
			name:    "successfully get transaction",
			uuid:    transaction.UUID,
			wantErr: false,
		},
		{
			name:    "get transaction not found",
			uuid:    "nonexistent-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetTransaction(ctx, tt.uuid)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, transaction.UUID, result.UUID)
			require.Equal(t, transaction.OrderUUID, result.OrderUUID)
		})
	}
}
