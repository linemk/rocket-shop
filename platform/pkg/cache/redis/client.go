package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/linemk/rocket-shop/platform/pkg/cache"
)

var ErrKeyNotFound = errors.New("key not found")

// client реализация cache.Client для Redis
type client struct {
	rdb         *redis.Client
	setOperator cache.SetOperator
}

// NewClient создает новый Redis клиент
func NewClient(cfg cache.Config) (cache.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
	})

	c := &client{
		rdb: rdb,
	}

	c.setOperator = NewSetOperator(rdb)

	return c, nil
}

func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func (c *client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrKeyNotFound
		}
		return nil, errors.Wrap(err, "failed to get value")
	}

	return result, nil
}

func (c *client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.Wrap(err, "failed to check key existence")
	}

	return result > 0, nil
}

func (c *client) SetOperator() cache.SetOperator {
	return c.setOperator
}
