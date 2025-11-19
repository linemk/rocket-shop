package service

import (
	"context"

	"github.com/linemk/rocket-shop/assembly/internal/model"
)

//go:generate ../../bin/mockery --name ConsumerService --output ./mocks --outpkg mocks --filename mock_consumer_service.go
type ConsumerService interface {
	RunConsumers(ctx context.Context) error
}

//go:generate ../../bin/mockery --name OrderProducerService --output ./mocks --outpkg mocks --filename mock_order_producer_service.go
type OrderProducerService interface {
	SendShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error
}
