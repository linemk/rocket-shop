package env

import (
	"fmt"
	"os"
)

const (
	paymentGRPCHostEnv = "PAYMENT_GRPC_HOST"
	paymentGRPCPortEnv = "PAYMENT_GRPC_PORT"
)

type paymentGRPCConfig struct {
	host string
	port string
}

// NewPaymentGRPCConfig создает конфигурацию gRPC сервера из переменных окружения
func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	host := os.Getenv(paymentGRPCHostEnv)
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv(paymentGRPCPortEnv)
	if port == "" {
		port = "50052"
	}

	return &paymentGRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (c *paymentGRPCConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
