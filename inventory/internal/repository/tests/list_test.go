package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	repoimpl "github.com/linemk/rocket-shop/inventory/internal/repository/inventory"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	ctx := context.Background()
	repo := repoimpl.NewRepository()

	// seed
	engine1 := models.Part{UUID: uuid.New().String(), Name: "Engine-1", Category: inventory_v1.Category_CATEGORY_ENGINE, Tags: []string{"hot", "metal"}}
	engine2 := models.Part{UUID: uuid.New().String(), Name: "Engine-2", Category: inventory_v1.Category_CATEGORY_ENGINE, Tags: []string{"metal"}}
	body := models.Part{UUID: uuid.New().String(), Name: "Body-1", Category: inventory_v1.Category_CATEGORY_WING, Tags: []string{"plastic"}}
	require.NoError(t, repo.CreatePart(ctx, engine1))
	require.NoError(t, repo.CreatePart(ctx, engine2))
	require.NoError(t, repo.CreatePart(ctx, body))

	t.Run("Empty filter returns all", func(t *testing.T) {
		got, err := repo.ListParts(ctx, models.PartFilter{})
		require.NoError(t, err)
		require.Len(t, got, 3)
	})

	t.Run("Filter by UUIDs", func(t *testing.T) {
		got, err := repo.ListParts(ctx, models.PartFilter{UUIDs: []string{engine1.UUID}})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, engine1.UUID, got[0].UUID)
	})

	t.Run("Filter by Names", func(t *testing.T) {
		got, err := repo.ListParts(ctx, models.PartFilter{Names: []string{"Engine-2"}})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, engine2.UUID, got[0].UUID)
	})

	t.Run("Filter by Categories", func(t *testing.T) {
		got, err := repo.ListParts(ctx, models.PartFilter{Categories: []inventory_v1.Category{inventory_v1.Category_CATEGORY_ENGINE}})
		require.NoError(t, err)
		require.Len(t, got, 2)
	})

	t.Run("Filter by Tags", func(t *testing.T) {
		got, err := repo.ListParts(ctx, models.PartFilter{Tags: []string{"plastic"}})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, body.UUID, got[0].UUID)
	})
}
