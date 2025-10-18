package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	"github.com/linemk/rocket-shop/inventory/internal/mocks"
	"github.com/linemk/rocket-shop/inventory/internal/usecase"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		repoMock func() *mocks.MockInventoryRepository
	}

	filter := models.PartFilter{
		Categories: []inventory_v1.Category{inventory_v1.Category_CATEGORY_ENGINE},
	}

	part1 := models.Part{
		UUID:          uuid.New().String(),
		Name:          "Engine Part 1",
		Price:         100.0,
		StockQuantity: 5,
		Category:      inventory_v1.Category_CATEGORY_ENGINE,
	}

	part2 := models.Part{
		UUID:          uuid.New().String(),
		Name:          "Engine Part 2",
		Price:         200.0,
		StockQuantity: 10,
		Category:      inventory_v1.Category_CATEGORY_ENGINE,
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "successfully list parts",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockRepo := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockRepo.EXPECT().ListParts(ctx, filter).Return([]models.Part{part1, part2}, nil)
					return mockRepo
				},
			},
			wantErr: false,
		},
		{
			name: "list parts with empty result",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockRepo := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockRepo.EXPECT().ListParts(ctx, filter).Return([]models.Part{}, nil)
					return mockRepo
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inventoryRepo := tt.fields.repoMock()
			uc := usecase.NewUseCase(inventoryRepo)

			parts, err := uc.ListParts(ctx, filter)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, parts)
		})
	}
}
