package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	v1 "github.com/linemk/rocket-shop/inventory/internal/delivery/v1"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	"github.com/linemk/rocket-shop/inventory/internal/mocks"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func TestListPartsAPI(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		useCaseMock func() *mocks.MockInventoryUseCase
	}

	part1 := models.PartInfo{
		UUID:          uuid.New().String(),
		Name:          "Engine Part 1",
		Price:         100.0,
		StockQuantity: 5,
		Category:      inventory_v1.Category_CATEGORY_ENGINE,
	}

	part2 := models.PartInfo{
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
			name: "successfully list parts via API",
			fields: fields{
				useCaseMock: func() *mocks.MockInventoryUseCase {
					mockUseCase := mocks.NewMockInventoryUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().ListParts(ctx, gomock.Any()).Return([]models.PartInfo{part1, part2}, nil)
					return mockUseCase
				},
			},
			wantErr: false,
		},
		{
			name: "successfully list parts with empty result",
			fields: fields{
				useCaseMock: func() *mocks.MockInventoryUseCase {
					mockUseCase := mocks.NewMockInventoryUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().ListParts(ctx, gomock.Any()).Return([]models.PartInfo{}, nil)
					return mockUseCase
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCaseMock := tt.fields.useCaseMock()
			api := v1.NewAPI(useCaseMock)

			req := &inventory_v1.ListPartsRequest{
				Filter: nil,
			}

			resp, err := api.ListParts(ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}
