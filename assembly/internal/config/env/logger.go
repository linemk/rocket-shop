package env

import (
	"os"
)

const (
	loggerLevelEnvName = "LOGGER_LEVEL"
)

type loggerConfig struct {
	level string
}

func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv(loggerLevelEnvName)
	if len(level) == 0 {
		level = "info"
	}

	return &loggerConfig{
		level: level,
	}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.level
}
