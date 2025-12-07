package env

import (
	"os"
	"time"
)

const (
	sessionTTLEnv = "SESSION_TTL"
)

type sessionConfig struct {
	ttl time.Duration
}

// NewSessionConfig создает конфигурацию сессий из переменных окружения
func NewSessionConfig() (*sessionConfig, error) {
	ttl := 24 * time.Hour // По умолчанию 24 часа

	if ttlStr := os.Getenv(sessionTTLEnv); ttlStr != "" {
		parsed, err := time.ParseDuration(ttlStr)
		if err == nil {
			ttl = parsed
		}
	}

	return &sessionConfig{
		ttl: ttl,
	}, nil
}

func (c *sessionConfig) TTL() time.Duration {
	return c.ttl
}
