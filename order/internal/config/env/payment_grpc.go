package env

import (
	"os"
)

const (
	paymentGRPCAddressEnv = "PAYMENT_GRPC_ADDRESS"
)

type paymentGRPCConfig struct {
	address string
}

// NewPaymentGRPCConfig создает конфигурацию gRPC клиента Payment из переменных окружения
func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	address := os.Getenv(paymentGRPCAddressEnv)
	if address == "" {
		address = "localhost:50052"
	}

	return &paymentGRPCConfig{
		address: address,
	}, nil
}

func (c *paymentGRPCConfig) Address() string {
	return c.address
}
