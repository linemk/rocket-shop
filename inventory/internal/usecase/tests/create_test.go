package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	"github.com/linemk/rocket-shop/inventory/internal/mocks"
	"github.com/linemk/rocket-shop/inventory/internal/usecase"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		repoMock func() *mocks.MockInventoryRepository
	}
	part := models.Part{}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockClient := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockClient.EXPECT().CreatePart(ctx, part).Return(nil)
					return mockClient
				},
			},
			wantErr: false,
		},
		{
			name: "Failure",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockClient := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockClient.EXPECT().CreatePart(ctx, part).Return(apperrors.ErrPartAlreadyExists)
					return mockClient
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inventoryRepo := tt.fields.repoMock()

			uc := usecase.NewUseCase(inventoryRepo)

			err := uc.CreatePart(ctx, part)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
