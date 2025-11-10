package env

import (
	"os"
)

const (
	logLevelEnv  = "LOG_LEVEL"
	logAsJSONEnv = "LOG_AS_JSON"
)

type loggerConfig struct {
	level  string
	asJSON bool
}

// NewLoggerConfig создает конфигурацию логгера из переменных окружения
func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv(logLevelEnv)
	if level == "" {
		level = "info"
	}

	asJSON := os.Getenv(logAsJSONEnv) == "true"

	return &loggerConfig{
		level:  level,
		asJSON: asJSON,
	}, nil
}

func (c *loggerConfig) Level() string {
	return c.level
}

func (c *loggerConfig) AsJSON() bool {
	return c.asJSON
}
