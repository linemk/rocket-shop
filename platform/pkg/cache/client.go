package cache

import (
	"context"
	"time"
)

// Config содержит конфигурацию для подключения к кэшу
type Config struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
}

// Client интерфейс для работы с кэшем
type Client interface {
	Ping(ctx context.Context) error
	Close() error

	// Базовые операции
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Set операции
	SetOperator() SetOperator
}

// SetOperator интерфейс для работы с множествами (Sets)
type SetOperator interface {
	// SAdd добавляет элементы в множество
	SAdd(ctx context.Context, key string, members ...string) error

	// SMembers возвращает все элементы множества
	SMembers(ctx context.Context, key string) ([]string, error)

	// SRem удаляет элементы из множества
	SRem(ctx context.Context, key string, members ...string) error

	// SCard возвращает количество элементов в множестве
	SCard(ctx context.Context, key string) (int64, error)
}
