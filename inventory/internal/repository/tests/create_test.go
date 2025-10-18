package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	repoimpl "github.com/linemk/rocket-shop/inventory/internal/repository/inventory"
)

func TestCreatePart(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		repo := repoimpl.NewRepository()

		part := models.Part{
			UUID:  uuid.New().String(),
			Name:  "Bolt X",
			Price: 9.99,
		}

		err := repo.CreatePart(ctx, part)
		require.NoError(t, err)

		got, err := repo.GetPart(ctx, part.UUID)
		require.NoError(t, err)
		require.Equal(t, part.UUID, got.UUID)
		require.Equal(t, part.Name, got.Name)
		require.Equal(t, part.Price, got.Price)
	})

	t.Run("Duplicate UUID", func(t *testing.T) {
		repo := repoimpl.NewRepository()

		id := uuid.New().String()
		part := models.Part{UUID: id, Name: "Washer", Price: 1.25}

		require.NoError(t, repo.CreatePart(ctx, part))

		err := repo.CreatePart(ctx, part)
		require.Error(t, err)
	})
}
