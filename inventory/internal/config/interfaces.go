package config

// LoggerConfig интерфейс конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// InventoryGRPCConfig интерфейс конфигурации gRPC сервера Inventory
type InventoryGRPCConfig interface {
	Address() string
}

// MongoConfig интерфейс конфигурации MongoDB
type MongoConfig interface {
	URI() string
	DatabaseName() string
}
