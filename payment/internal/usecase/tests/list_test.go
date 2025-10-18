package tests

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
	"github.com/linemk/rocket-shop/payment/internal/mocks"
	"github.com/linemk/rocket-shop/payment/internal/usecase"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func TestListTransactions(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		repoMock func() *mocks.MockPaymentRepository
	}

	orderUUID := uuid.New().String()

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
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_SBP,
		Amount:        25.50,
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
			name: "successfully list transactions",
			fields: fields{
				repoMock: func() *mocks.MockPaymentRepository {
					mockRepo := mocks.NewMockPaymentRepository(gomock.NewController(t))
					mockRepo.EXPECT().ListTransactions(ctx, orderUUID).Return([]models.Transaction{transaction1, transaction2}, nil)
					return mockRepo
				},
			},
			uuid:    orderUUID,
			wantErr: false,
		},
		{
			name: "list transactions with empty result",
			fields: fields{
				repoMock: func() *mocks.MockPaymentRepository {
					mockRepo := mocks.NewMockPaymentRepository(gomock.NewController(t))
					mockRepo.EXPECT().ListTransactions(ctx, orderUUID).Return([]models.Transaction{}, nil)
					return mockRepo
				},
			},
			uuid:    orderUUID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepo := tt.fields.repoMock()
			uc := usecase.NewUseCase(paymentRepo)

			transactions, err := uc.ListTransactions(ctx, tt.uuid)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, transactions)
		})
	}
}
