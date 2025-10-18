package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	"github.com/linemk/rocket-shop/inventory/internal/mocks"
	"github.com/linemk/rocket-shop/inventory/internal/usecase"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func TestUpdatePart(t *testing.T) {
	ctx := context.Background()
	testUUID := uuid.New().String()

	type fields struct {
		repoMock func() *mocks.MockInventoryRepository
	}

	updatePart := models.Part{
		UUID:          testUUID,
		Name:          "Updated Engine Part",
		Description:   "Updated description",
		Price:         150.0,
		StockQuantity: 10,
		Category:      inventory_v1.Category_CATEGORY_ENGINE,
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "successfully update a part",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockRepo := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockRepo.EXPECT().UpdatePart(ctx, testUUID, updatePart).Return(nil)
					return mockRepo
				},
			},
			wantErr: false,
		},
		{
			name: "error part not found",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockRepo := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockRepo.EXPECT().UpdatePart(ctx, testUUID, updatePart).Return(apperrors.ErrPartNotFound)
					return mockRepo
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inventoryRepo := tt.fields.repoMock()
			uc := usecase.NewUseCase(inventoryRepo)

			err := uc.UpdatePart(ctx, testUUID, updatePart)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
