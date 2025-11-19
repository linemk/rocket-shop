package order_assembled_consumer

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
	h.logger.Info(ctx, "Received OrderAssembled event", zap.String("topic", msg.Topic))

	// Декодируем событие
	event, err := decoder.DecodeShipAssembled(msg.Value)
	if err != nil {
		h.logger.Error(ctx, "Failed to decode OrderAssembled event", zap.Error(err))
		return err
	}

	h.logger.Info(ctx, "Processing order assembled notification",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	// Отправляем уведомление в Telegram
	if err := h.telegramService.SendOrderAssembledNotification(ctx, event); err != nil {
		h.logger.Error(ctx, "Failed to send order assembled notification", zap.Error(err))
		return err
	}

	h.logger.Info(ctx, "Order assembled notification processed successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
