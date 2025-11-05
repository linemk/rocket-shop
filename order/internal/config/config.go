package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/linemk/rocket-shop/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	OrderHTTP     OrderHTTPConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
}

// Load загружает конфигурацию из переменных окружения
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		OrderHTTP:     orderHTTPCfg,
		Postgres:      postgresCfg,
		InventoryGRPC: inventoryGRPCCfg,
		PaymentGRPC:   paymentGRPCCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию приложения
func AppConfig() *config {
	return appConfig
}
