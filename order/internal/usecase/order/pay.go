package order

import (
	"context"
	"fmt"

	"github.com/linemk/rocket-shop/order/internal/client/grpc/payment/converter"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (uc *useCase) PayOrder(ctx context.Context, uuid string, paymentMethod order_v1.PaymentMethod) (string, error) {
	// 1. Получаем заказ
	order, err := uc.orderRepository.Get(ctx, uuid)
	if err != nil {
		return "", fmt.Errorf("order not found: %w", err)
	}

	// 2. Бизнес-логика: проверяем статус
	if order.Status != order_v1.OrderStatusPENDINGPAYMENT {
		return "", fmt.Errorf("order cannot be paid in current status")
	}

	// 3. Вызываем PaymentService
	protoPaymentMethod := converter.OpenAPIPaymentMethodToProto(paymentMethod)
	transactionUUID, err := uc.paymentClient.PayOrder(ctx, order.UUID, order.UserID, protoPaymentMethod)
	if err != nil {
		return "", fmt.Errorf("payment failed: %w", err)
	}

	// 4. Обновляем заказ
	updateInfo := models.OrderUpdateInfo{
		Status:        &[]order_v1.OrderStatus{order_v1.OrderStatusPAID}[0],
		TransactionID: &transactionUUID,
		PaymentMethod: &paymentMethod,
	}

	if err := uc.orderRepository.Update(ctx, uuid, updateInfo); err != nil {
		return "", fmt.Errorf("failed to update order: %w", err)
	}

	return transactionUUID, nil
}
