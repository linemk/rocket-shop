package config

// LoggerConfig интерфейс конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// OrderHTTPConfig интерфейс конфигурации HTTP сервера Order
type OrderHTTPConfig interface {
	Address() string
}

// PostgresConfig интерфейс конфигурации PostgreSQL
type PostgresConfig interface {
	DSN() string
	MigrationsDir() string
}

// InventoryGRPCConfig интерфейс конфигурации gRPC клиента для Inventory
type InventoryGRPCConfig interface {
	Address() string
}

// PaymentGRPCConfig интерфейс конфигурации gRPC клиента для Payment
type PaymentGRPCConfig interface {
	Address() string
}
