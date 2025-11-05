package env

import (
	"os"
)

const (
	inventoryGRPCAddressEnv = "INVENTORY_GRPC_ADDRESS"
)

type inventoryGRPCConfig struct {
	address string
}

// NewInventoryGRPCConfig создает конфигурацию gRPC клиента Inventory из переменных окружения
func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	address := os.Getenv(inventoryGRPCAddressEnv)
	if address == "" {
		address = "localhost:50051"
	}

	return &inventoryGRPCConfig{
		address: address,
	}, nil
}

func (c *inventoryGRPCConfig) Address() string {
	return c.address
}
