package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/linemk/rocket-shop/order/internal/converter/kafka/decoder"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	"github.com/linemk/rocket-shop/order/internal/repository"
	platformKafka "github.com/linemk/rocket-shop/platform/pkg/kafka"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

type handler struct {
	orderRepository repository.OrderRepository
	logger          Logger
}

func NewHandler(orderRepository repository.OrderRepository, logger Logger) HandlerFunc {
	h := &handler{
		orderRepository: orderRepository,
		logger:          logger,
	}

	return h.Handle
}

func (h *handler) Handle(ctx context.Context, msg platformKafka.Message) error {
	h.logger.Info(ctx, "Received ShipAssembled event", zap.String("topic", msg.Topic))

	// Декодируем событие
	event, err := decoder.DecodeShipAssembled(msg.Value)
	if err != nil {
		h.logger.Error(ctx, "Failed to decode ShipAssembled event", zap.Error(err))
		return err
	}

	h.logger.Info(ctx, "Processing ship assembled event",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	// Обновляем статус заказа на COMPLETED
	updateInfo := models.OrderUpdateInfo{
		Status: &[]order_v1.OrderStatus{order_v1.OrderStatusCOMPLETED}[0],
	}

	if err := h.orderRepository.Update(ctx, event.OrderUUID, updateInfo); err != nil {
		h.logger.Error(ctx, "Failed to update order status to ASSEMBLED", zap.Error(err))
		return err
	}

	h.logger.Info(ctx, "Order status updated to ASSEMBLED successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
