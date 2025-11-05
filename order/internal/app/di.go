package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	inventoryClient "github.com/linemk/rocket-shop/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/linemk/rocket-shop/order/internal/client/grpc/payment/v1"
	"github.com/linemk/rocket-shop/order/internal/config"
	v1 "github.com/linemk/rocket-shop/order/internal/delivery/v1"
	"github.com/linemk/rocket-shop/order/internal/repository"
	"github.com/linemk/rocket-shop/order/internal/usecase"
	"github.com/linemk/rocket-shop/platform/pkg/closer"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

type diContainer struct {
	orderV1API order_v1.Handler

	orderUseCase usecase.OrderUseCase

	orderRepository repository.OrderRepository

	inventoryClient inventoryClient.InventoryClient
	paymentClient   paymentClient.PaymentClient

	dbPool *pgxpool.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) SetDBPool(pool *pgxpool.Pool) {
	d.dbPool = pool
}

func (d *diContainer) OrderV1API(ctx context.Context) order_v1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = v1.NewAPI(d.OrderUseCase(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderUseCase(ctx context.Context) usecase.OrderUseCase {
	if d.orderUseCase == nil {
		d.orderUseCase = usecase.NewUseCase(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderUseCase
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = repository.NewRepository(d.dbPool)
	}

	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) inventoryClient.InventoryClient {
	if d.inventoryClient == nil {
		client, err := inventoryClient.NewClient(config.AppConfig().InventoryGRPC.Address())
		if err != nil {
			panic(fmt.Sprintf("failed to create inventory client: %s\n", err.Error()))
		}

		closer.AddNamed("Inventory gRPC client", func(ctx context.Context) error {
			return client.Close()
		})

		d.inventoryClient = client
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) paymentClient.PaymentClient {
	if d.paymentClient == nil {
		client, err := paymentClient.NewClient(config.AppConfig().PaymentGRPC.Address())
		if err != nil {
			panic(fmt.Sprintf("failed to create payment client: %s\n", err.Error()))
		}

		closer.AddNamed("Payment gRPC client", func(ctx context.Context) error {
			return client.Close()
		})

		d.paymentClient = client
	}

	return d.paymentClient
}
