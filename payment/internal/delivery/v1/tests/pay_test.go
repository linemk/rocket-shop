package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	v1 "github.com/linemk/rocket-shop/payment/internal/delivery/v1"
	"github.com/linemk/rocket-shop/payment/internal/mocks"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func TestPayOrderAPI(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		useCaseMock func() *mocks.MockPaymentUseCase
	}

	type args struct {
		req *payment_v1.PayOrderRequest
	}

	orderUUID := uuid.New().String()
	userUUID := uuid.New().String()
	transactionUUID := uuid.New().String()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successfully pay order via API",
			fields: fields{
				useCaseMock: func() *mocks.MockPaymentUseCase {
					mockUseCase := mocks.NewMockPaymentUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().PayOrder(ctx, orderUUID, userUUID, payment_v1.PaymentMethod_PAYMENT_METHOD_CARD).Return(transactionUUID, nil)
					return mockUseCase
				},
			},
			args: args{
				req: &payment_v1.PayOrderRequest{
					OrderUuid:     orderUUID,
					UserUuid:      userUUID,
					PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			wantErr: false,
		},
		{
			name: "pay order with empty orderUUID",
			fields: fields{
				useCaseMock: func() *mocks.MockPaymentUseCase {
					mockUseCase := mocks.NewMockPaymentUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().PayOrder(ctx, "", userUUID, gomock.Any()).Return("", fmt.Errorf("invalid order UUID")).AnyTimes()
					return mockUseCase
				},
			},
			args: args{
				req: &payment_v1.PayOrderRequest{
					OrderUuid:     "",
					UserUuid:      userUUID,
					PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			wantErr: true,
		},
		{
			name: "pay order with empty userUUID",
			fields: fields{
				useCaseMock: func() *mocks.MockPaymentUseCase {
					mockUseCase := mocks.NewMockPaymentUseCase(gomock.NewController(t))
					mockUseCase.EXPECT().PayOrder(ctx, orderUUID, "", gomock.Any()).Return("", fmt.Errorf("invalid user UUID")).AnyTimes()
					return mockUseCase
				},
			},
			args: args{
				req: &payment_v1.PayOrderRequest{
					OrderUuid:     orderUUID,
					UserUuid:      "",
					PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCaseMock := tt.fields.useCaseMock()
			api := v1.NewAPI(useCaseMock)

			resp, err := api.PayOrder(ctx, tt.args.req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, transactionUUID, resp.TransactionUuid)
		})
	}
}
