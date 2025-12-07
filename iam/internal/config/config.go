package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/linemk/rocket-shop/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	GRPC     GRPCConfig
	Session  SessionConfig
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

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerCfg,
		Postgres: postgresCfg,
		Redis:    redisCfg,
		GRPC:     grpcCfg,
		Session:  sessionCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию приложения
func AppConfig() *config {
	return appConfig
}
