package env

import (
	"fmt"
	"os"
)

const (
	inventoryGRPCHostEnv = "INVENTORY_GRPC_HOST"
	inventoryGRPCPortEnv = "INVENTORY_GRPC_PORT"
)

type inventoryGRPCConfig struct {
	host string
	port string
}

// NewInventoryGRPCConfig создает конфигурацию gRPC сервера из переменных окружения
func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	host := os.Getenv(inventoryGRPCHostEnv)
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv(inventoryGRPCPortEnv)
	if port == "" {
		port = "50051"
	}

	return &inventoryGRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (c *inventoryGRPCConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
