package env

import (
	"fmt"
	"os"
)

const (
	orderHTTPHostEnv = "ORDER_HTTP_HOST"
	orderHTTPPortEnv = "ORDER_HTTP_PORT"
)

type orderHTTPConfig struct {
	host string
	port string
}

// NewOrderHTTPConfig создает конфигурацию HTTP сервера из переменных окружения
func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	host := os.Getenv(orderHTTPHostEnv)
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv(orderHTTPPortEnv)
	if port == "" {
		port = "8080"
	}

	return &orderHTTPConfig{
		host: host,
		port: port,
	}, nil
}

func (c *orderHTTPConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
