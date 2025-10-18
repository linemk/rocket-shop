package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	v1 "github.com/linemk/rocket-shop/order/internal/delivery/v1"
	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	"github.com/linemk/rocket-shop/order/internal/mocks"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	orderUUID := uuid.New()

	type fields struct {
		orderUseCase func() *mocks.MockOrderUseCase
	}

	tests := []struct {
		name    string
		fields  fields
		req     order_v1.OptCreateOrderReq
		wantErr bool
	}{
		{
			name: "successful create order",
			fields: fields{
				orderUseCase: func() *mocks.MockOrderUseCase {
					mockClient := mocks.NewMockOrderUseCase(gomock.NewController(t))
					mockClient.EXPECT().CreateOrder(ctx, gomock.Any()).Return(orderUUID.String(), nil)
					mockClient.EXPECT().GetOrder(ctx, orderUUID.String()).Return(models.Order{
						UUID:       orderUUID.String(),
						TotalPrice: 300.0,
					}, nil)

					return mockClient
				},
			},
			req: order_v1.OptCreateOrderReq{
				Value: order_v1.CreateOrderReq{
					UserUUID:  uuid.New(),
					PartUuids: []uuid.UUID{partUUID1, partUUID2},
				},
			},
			wantErr: false,
		},
		{
			name: "error create order",
			fields: fields{
				orderUseCase: func() *mocks.MockOrderUseCase {
					mockClient := mocks.NewMockOrderUseCase(gomock.NewController(t))
					mockClient.EXPECT().CreateOrder(ctx, gomock.Any()).Return("", apperrors.ErrNoPartsSpecified)

					return mockClient
				},
			},
			req: order_v1.OptCreateOrderReq{
				Value: order_v1.CreateOrderReq{
					UserUUID:  uuid.New(),
					PartUuids: []uuid.UUID{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderUseCase := tt.fields.orderUseCase()

			api := v1.NewAPI(orderUseCase)

			result, err := api.CreateOrder(ctx, tt.req)

			if tt.wantErr {
				// Для невалидного запроса API возвращает BadRequest в теле ответа и err == nil
				require.NoError(t, err)
				require.IsType(t, &order_v1.BadRequest{}, result)
				bad := result.(*order_v1.BadRequest)
				require.Equal(t, 400, bad.Code)
				return
			}

			require.NoError(t, err)
			require.IsType(t, &order_v1.CreateOrderResp{}, result)
			resp := result.(*order_v1.CreateOrderResp)
			require.Equal(t, orderUUID, resp.UUID)
			require.Equal(t, float32(300.0), resp.TotalPrice)
		})
	}
}
