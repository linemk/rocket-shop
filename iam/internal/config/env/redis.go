package env

import (
	"os"
	"strconv"
	"time"
)

const (
	redisAddrEnv         = "REDIS_ADDR"
	redisPasswordEnv     = "REDIS_PASSWORD"
	redisDBEnv           = "REDIS_DB"
	redisDialTimeoutEnv  = "REDIS_DIAL_TIMEOUT"
	redisReadTimeoutEnv  = "REDIS_READ_TIMEOUT"
	redisWriteTimeoutEnv = "REDIS_WRITE_TIMEOUT"
	redisPoolSizeEnv     = "REDIS_POOL_SIZE"
)

type redisConfig struct {
	addr         string
	password     string
	db           int
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	poolSize     int
}

// NewRedisConfig создает конфигурацию Redis из переменных окружения
func NewRedisConfig() (*redisConfig, error) {
	addr := os.Getenv(redisAddrEnv)
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv(redisPasswordEnv)

	db := 0
	if dbStr := os.Getenv(redisDBEnv); dbStr != "" {
		parsed, err := strconv.Atoi(dbStr)
		if err == nil {
			db = parsed
		}
	}

	dialTimeout := 5 * time.Second
	if timeoutStr := os.Getenv(redisDialTimeoutEnv); timeoutStr != "" {
		parsed, err := time.ParseDuration(timeoutStr)
		if err == nil {
			dialTimeout = parsed
		}
	}

	readTimeout := 3 * time.Second
	if timeoutStr := os.Getenv(redisReadTimeoutEnv); timeoutStr != "" {
		parsed, err := time.ParseDuration(timeoutStr)
		if err == nil {
			readTimeout = parsed
		}
	}

	writeTimeout := 3 * time.Second
	if timeoutStr := os.Getenv(redisWriteTimeoutEnv); timeoutStr != "" {
		parsed, err := time.ParseDuration(timeoutStr)
		if err == nil {
			writeTimeout = parsed
		}
	}

	poolSize := 10
	if poolSizeStr := os.Getenv(redisPoolSizeEnv); poolSizeStr != "" {
		parsed, err := strconv.Atoi(poolSizeStr)
		if err == nil {
			poolSize = parsed
		}
	}

	return &redisConfig{
		addr:         addr,
		password:     password,
		db:           db,
		dialTimeout:  dialTimeout,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		poolSize:     poolSize,
	}, nil
}

func (c *redisConfig) Addr() string {
	return c.addr
}

func (c *redisConfig) Password() string {
	return c.password
}

func (c *redisConfig) DB() int {
	return c.db
}

func (c *redisConfig) DialTimeout() time.Duration {
	return c.dialTimeout
}

func (c *redisConfig) ReadTimeout() time.Duration {
	return c.readTimeout
}

func (c *redisConfig) WriteTimeout() time.Duration {
	return c.writeTimeout
}

func (c *redisConfig) PoolSize() int {
	return c.poolSize
}
