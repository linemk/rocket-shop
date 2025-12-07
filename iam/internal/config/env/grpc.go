package env

import (
	"os"
)

const (
	grpcAddressEnv = "IAM_GRPC_ADDRESS"
)

type grpcConfig struct {
	address string
}

// NewGRPCConfig создает конфигурацию gRPC сервера из переменных окружения
func NewGRPCConfig() (*grpcConfig, error) {
	address := os.Getenv(grpcAddressEnv)
	if address == "" {
		address = ":50053"
	}

	return &grpcConfig{
		address: address,
	}, nil
}

func (c *grpcConfig) Address() string {
	return c.address
}
