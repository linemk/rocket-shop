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

func TestGetPart(t *testing.T) {
	ctx := context.Background()
	testUUID := uuid.New().String()

	type fields struct {
		repoMock func() *mocks.MockInventoryRepository
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "successfully get a part",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockRepo := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockRepo.EXPECT().GetPart(ctx, testUUID).Return(
						models.Part{
							UUID:          testUUID,
							Name:          "Engine Part",
							Description:   "High performance engine component",
							Price:         100.0,
							StockQuantity: 5,
							Category:      inventory_v1.Category_CATEGORY_ENGINE,
						}, nil,
					)
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
					mockRepo.EXPECT().GetPart(ctx, testUUID).Return(models.Part{}, apperrors.ErrPartNotFound)
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

			partInfo, err := uc.GetPart(ctx, testUUID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, testUUID, partInfo.UUID)
			require.Equal(t, "Engine Part", partInfo.Name)
		})
	}
}
