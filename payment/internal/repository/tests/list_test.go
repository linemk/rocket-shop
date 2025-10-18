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

func TestListTransactions(t *testing.T) {
	ctx := context.Background()

	repo := payment.NewRepository()

	orderUUID := uuid.New().String()
	otherOrderUUID := uuid.New().String()

	// Create transactions for first order
	transaction1 := models.Transaction{
		UUID:          uuid.New().String(),
		OrderUUID:     orderUUID,
		UserID:        uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
		Amount:        50.00,
		Status:        models.TransactionStatusCompleted,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	transaction2 := models.Transaction{
		UUID:          uuid.New().String(),
		OrderUUID:     orderUUID,
		UserID:        uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
		Amount:        25.50,
		Status:        models.TransactionStatusCompleted,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create transaction for different order
	transaction3 := models.Transaction{
		UUID:          uuid.New().String(),
		OrderUUID:     otherOrderUUID,
		UserID:        uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
		Amount:        100.00,
		Status:        models.TransactionStatusCompleted,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Populate repo
	err := repo.CreateTransaction(ctx, transaction1)
	require.NoError(t, err)
	err = repo.CreateTransaction(ctx, transaction2)
	require.NoError(t, err)
	err = repo.CreateTransaction(ctx, transaction3)
	require.NoError(t, err)

	tests := []struct {
		name        string
		orderUUID   string
		expectedLen int
		wantErr     bool
	}{
		{
			name:        "successfully list transactions for order",
			orderUUID:   orderUUID,
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name:        "list transactions for different order",
			orderUUID:   otherOrderUUID,
			expectedLen: 1,
			wantErr:     false,
		},
		{
			name:        "list transactions for nonexistent order",
			orderUUID:   uuid.New().String(),
			expectedLen: 0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ListTransactions(ctx, tt.orderUUID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, result, tt.expectedLen)
		})
	}
}
