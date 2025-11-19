package service

import (
	"context"

	"github.com/linemk/rocket-shop/order/internal/entyties/events"
)

//go:generate ../../bin/mockery --name ConsumerService --output ../mocks --outpkg mocks --filename consumer_service_mock.go
type ConsumerService interface {
	RunConsumers(ctx context.Context) error
}

//go:generate ../../bin/mockery --name OrderProducerService --output ../mocks --outpkg mocks --filename order_producer_service_mock.go
type OrderProducerService interface {
	SendOrderPaid(ctx context.Context, event *events.OrderPaidEvent) error
}
