package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	v1 "github.com/linemk/rocket-shop/inventory/internal/delivery/v1"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	"github.com/linemk/rocket-shop/inventory/internal/mocks"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func TestGetPartAPI(t *testing.T) {
	ctx := context.Background()
	testUUID := uuid.New().String()

	type fields struct {
		useCaseMock func() *mocks.MockInventoryUseCase
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "successfully get part via API",
			fields: fields{
				useCaseMock: func() *mocks.MockInventoryUseCase {
					mockUseCase := mocks.NewMockInventoryUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().GetPart(ctx, testUUID).Return(
						models.PartInfo{
							UUID:          testUUID,
							Name:          "Engine Part",
							Price:         100.0,
							StockQuantity: 5,
							Category:      inventory_v1.Category_CATEGORY_ENGINE,
						}, nil,
					)
					return mockUseCase
				},
			},
			wantErr: false,
		},
		{
			name: "error part not found",
			fields: fields{
				useCaseMock: func() *mocks.MockInventoryUseCase {
					mockUseCase := mocks.NewMockInventoryUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().GetPart(ctx, testUUID).Return(models.PartInfo{}, apperrors.ErrPartNotFound)
					return mockUseCase
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCaseMock := tt.fields.useCaseMock()
			api := v1.NewAPI(useCaseMock)

			req := &inventory_v1.GetPartRequest{
				Uuid: testUUID,
			}

			resp, err := api.GetPart(ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, testUUID, resp.Part.Uuid)
		})
	}
}
