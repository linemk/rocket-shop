package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/mocks"
	"github.com/linemk/rocket-shop/inventory/internal/usecase"
)

func TestDeletePart(t *testing.T) {
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
			name: "successfully delete a part",
			fields: fields{
				repoMock: func() *mocks.MockInventoryRepository {
					mockRepo := mocks.NewMockInventoryRepository(gomock.NewController(t))
					mockRepo.EXPECT().DeletePart(ctx, testUUID).Return(nil)
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
					mockRepo.EXPECT().DeletePart(ctx, testUUID).Return(apperrors.ErrPartNotFound)
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

			err := uc.DeletePart(ctx, testUUID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
