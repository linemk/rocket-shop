package telegram

import (
	"bytes"
	"context"
	_ "embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/linemk/rocket-shop/notification/internal/model"
)

//go:embed templates/paid_notification.tmpl
var paidNotificationTemplate string

//go:embed templates/assembled_notification.tmpl
var assembledNotificationTemplate string

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type TelegramClient interface {
	SendMessage(ctx context.Context, text string) error
}

type service struct {
	client TelegramClient
	logger Logger
}

func NewService(client TelegramClient, logger Logger) *service {
	return &service{
		client: client,
		logger: logger,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event *model.OrderPaidEvent) error {
	s.logger.Info(ctx, "Sending order paid notification",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	message, err := s.renderPaidNotification(event)
	if err != nil {
		s.logger.Error(ctx, "Failed to render paid notification", zap.Error(err))
		return err
	}

	if err := s.client.SendMessage(ctx, message); err != nil {
		s.logger.Error(ctx, "Failed to send paid notification", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Order paid notification sent successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, event *model.ShipAssembledEvent) error {
	s.logger.Info(ctx, "Sending order assembled notification",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	message, err := s.renderAssembledNotification(event)
	if err != nil {
		s.logger.Error(ctx, "Failed to render assembled notification", zap.Error(err))
		return err
	}

	if err := s.client.SendMessage(ctx, message); err != nil {
		s.logger.Error(ctx, "Failed to send assembled notification", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Order assembled notification sent successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}

func (s *service) renderPaidNotification(event *model.OrderPaidEvent) (string, error) {
	tmpl, err := template.New("paid").Parse(paidNotificationTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, event); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) renderAssembledNotification(event *model.ShipAssembledEvent) (string, error) {
	tmpl, err := template.New("assembled").Parse(assembledNotificationTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, event); err != nil {
		return "", err
	}

	return buf.String(), nil
}
