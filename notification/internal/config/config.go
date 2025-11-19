package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/linemk/rocket-shop/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	TelegramBot            TelegramBotConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	telegramBotCfg, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		TelegramBot:            telegramBotCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию приложения
func AppConfig() *config {
	return appConfig
}
