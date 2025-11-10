package config

// LoggerConfig интерфейс конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// PaymentGRPCConfig интерфейс конфигурации gRPC сервера Payment
type PaymentGRPCConfig interface {
	Address() string
}
