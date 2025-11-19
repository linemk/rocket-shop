package service

import (
	"context"

	"github.com/linemk/rocket-shop/notification/internal/model"
)

//go:generate ../../bin/mockery --name ConsumerService --output ./mocks --outpkg mocks --filename mock_consumer_service.go
type ConsumerService interface {
	RunConsumers(ctx context.Context) error
}

//go:generate ../../bin/mockery --name TelegramService --output ./mocks --outpkg mocks --filename mock_telegram_service.go
type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event *model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, event *model.ShipAssembledEvent) error
}
