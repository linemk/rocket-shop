package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/linemk/rocket-shop/notification/internal/converter/kafka/decoder"
	"github.com/linemk/rocket-shop/notification/internal/service"
	platformKafka "github.com/linemk/rocket-shop/platform/pkg/kafka"
)

type handler struct {
	telegramService service.TelegramService
	logger          Logger
}

func NewHandler(telegramService service.TelegramService, logger Logger) platformKafka.MessageHandler {
	h := &handler{
		telegramService: telegramService,
		logger:          logger,
	}

	return h.Handle
}

func (h *handler) Handle(ctx context.Context, msg platformKafka.Message) error {
	h.logger.Info(ctx, "Received OrderPaid event", zap.String("topic", msg.Topic))

	// Декодируем событие
	event, err := decoder.DecodeOrderPaid(msg.Value)
	if err != nil {
		h.logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return err
	}

	h.logger.Info(ctx, "Processing order paid notification",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	// Отправляем уведомление в Telegram
	if err := h.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		h.logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
		return err
	}

	h.logger.Info(ctx, "Order paid notification processed successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
