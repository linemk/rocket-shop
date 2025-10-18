package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	repoimpl "github.com/linemk/rocket-shop/inventory/internal/repository/inventory"
)

func TestGetPart(t *testing.T) {
	ctx := context.Background()

	t.Run("Found", func(t *testing.T) {
		repo := repoimpl.NewRepository()
		id := uuid.New().String()
		want := models.Part{UUID: id, Name: "Nut", Price: 0.99}
		require.NoError(t, repo.CreatePart(ctx, want))

		got, err := repo.GetPart(ctx, id)
		require.NoError(t, err)
		require.Equal(t, want.UUID, got.UUID)
		require.Equal(t, want.Name, got.Name)
		require.Equal(t, want.Price, got.Price)
	})

	t.Run("Not found", func(t *testing.T) {
		repo := repoimpl.NewRepository()
		_, err := repo.GetPart(ctx, uuid.New().String())
		require.Error(t, err)
	})
}
