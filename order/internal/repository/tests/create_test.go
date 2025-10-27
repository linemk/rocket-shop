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

	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	"github.com/linemk/rocket-shop/order/internal/repository"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func TestCreate(t *testing.T) {
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
		PartUUIDs:     []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:    300.0,
		TransactionID: "",
		PaymentMethod: order_v1.PaymentMethodPAYMENTMETHODCARD,
		Status:        order_v1.OrderStatusPENDINGPAYMENT,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Проверяем что заказ создался
	retrievedOrder, err := repo.Get(ctx, orderUUID.String())
	require.NoError(t, err)
	require.Equal(t, order.UUID, retrievedOrder.UUID)
	require.Equal(t, order.UserID, retrievedOrder.UserID)
	require.Equal(t, order.TotalPrice, retrievedOrder.TotalPrice)
}
