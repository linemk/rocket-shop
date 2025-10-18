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

func TestCreateTransaction(t *testing.T) {
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

	tests := []struct {
		name    string
		trans   models.Transaction
		wantErr bool
	}{
		{
			name:    "successfully create transaction",
			trans:   transaction,
			wantErr: false,
		},
		{
			name:    "create transaction with duplicate UUID",
			trans:   transaction, // Same UUID
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTransaction(ctx, tt.trans)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
