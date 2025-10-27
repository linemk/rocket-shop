//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	"github.com/linemk/rocket-shop/order/internal/repository"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	// This test requires a real PostgreSQL database
	// Run with: go test -tags=integration ./...
	pool, err := pgxpool.New(ctx, "postgres://order_user:order_password@localhost:5432/order_db?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	repo := repository.NewRepository(pool)

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
		name       string
		uuid       string
		updateInfo models.OrderUpdateInfo
		wantErr    bool
	}{
		{
			name: "successful update order status",
			uuid: orderUUID.String(),
			updateInfo: models.OrderUpdateInfo{
				Status: &[]order_v1.OrderStatus{order_v1.OrderStatusPAID}[0],
			},
			wantErr: false,
		},
		{
			name: "successful update transaction id",
			uuid: orderUUID.String(),
			updateInfo: models.OrderUpdateInfo{
				TransactionID: &[]string{"tx-123"}[0],
			},
			wantErr: false,
		},
		{
			name: "error order not found",
			uuid: uuid.New().String(),
			updateInfo: models.OrderUpdateInfo{
				Status: &[]order_v1.OrderStatus{order_v1.OrderStatusPAID}[0],
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.uuid, tt.updateInfo)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, apperrors.ErrOrderNotFound, err)
				return
			}

			require.NoError(t, err)

			// Проверяем что обновление прошло
			updatedOrder, err := repo.Get(ctx, tt.uuid)
			require.NoError(t, err)
			require.NotNil(t, updatedOrder.UpdatedAt)

			if tt.updateInfo.Status != nil {
				require.Equal(t, *tt.updateInfo.Status, updatedOrder.Status)
			}
			if tt.updateInfo.TransactionID != nil {
				require.Equal(t, *tt.updateInfo.TransactionID, updatedOrder.TransactionID)
			}
		})
	}
}
