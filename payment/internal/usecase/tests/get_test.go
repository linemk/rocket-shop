package tests

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/payment/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
	"github.com/linemk/rocket-shop/payment/internal/mocks"
	"github.com/linemk/rocket-shop/payment/internal/usecase"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func TestGetTransaction(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		repoMock func() *mocks.MockPaymentRepository
	}

	transactionUUID := uuid.New().String()
	transaction := models.Transaction{
		UUID:          transactionUUID,
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
		fields  fields
		uuid    string
		wantErr bool
	}{
		{
			name: "successfully get transaction",
			fields: fields{
				repoMock: func() *mocks.MockPaymentRepository {
					mockRepo := mocks.NewMockPaymentRepository(gomock.NewController(t))
					mockRepo.EXPECT().GetTransaction(ctx, transactionUUID).Return(transaction, nil)
					return mockRepo
				},
			},
			uuid:    transactionUUID,
			wantErr: false,
		},
		{
			name: "get transaction not found",
			fields: fields{
				repoMock: func() *mocks.MockPaymentRepository {
					mockRepo := mocks.NewMockPaymentRepository(gomock.NewController(t))
					mockRepo.EXPECT().GetTransaction(ctx, gomock.Any()).Return(models.Transaction{}, apperrors.ErrTransactionNotFound)
					return mockRepo
				},
			},
			uuid:    "nonexistent-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepo := tt.fields.repoMock()
			uc := usecase.NewUseCase(paymentRepo)

			result, err := uc.GetTransaction(ctx, tt.uuid)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, transaction.UUID, result.UUID)
		})
	}
}
