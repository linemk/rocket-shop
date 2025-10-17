package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	"github.com/linemk/rocket-shop/order/internal/repository"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewRepository()

	orderUUID := uuid.New()
	now := time.Now()

	order := models.Order{
		UUID:          orderUUID.String(),
		UserID:        "user-123",
		PartUUIDs:     []uuid.UUID{uuid.New()},
		TotalPrice:    100.0,
		TransactionID: "",
		PaymentMethod: order_v1.PaymentMethodPAYMENTMETHODCARD,
		Status:        order_v1.OrderStatusPENDINGPAYMENT,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	// Создаём заказ
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{
			name:    "successful get order",
			uuid:    orderUUID.String(),
			wantErr: false,
		},
		{
			name:    "error order not found",
			uuid:    uuid.New().String(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Get(ctx, tt.uuid)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, apperrors.ErrOrderNotFound, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, order.UUID, result.UUID)
			require.Equal(t, order.UserID, result.UserID)
		})
	}
}
