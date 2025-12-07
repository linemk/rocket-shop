package redis

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/linemk/rocket-shop/platform/pkg/cache"
)

// setOperator реализация cache.SetOperator для Redis
type setOperator struct {
	rdb *redis.Client
}

// NewSetOperator создает новый SetOperator
func NewSetOperator(rdb *redis.Client) cache.SetOperator {
	return &setOperator{
		rdb: rdb,
	}
}

func (s *setOperator) SAdd(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	// Конвертируем []string в []interface{} для Redis
	interfaceMembers := make([]interface{}, len(members))
	for i, member := range members {
		interfaceMembers[i] = member
	}

	return s.rdb.SAdd(ctx, key, interfaceMembers...).Err()
}

func (s *setOperator) SMembers(ctx context.Context, key string) ([]string, error) {
	result, err := s.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get set members")
	}

	return result, nil
}

func (s *setOperator) SRem(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	// Конвертируем []string в []interface{} для Redis
	interfaceMembers := make([]interface{}, len(members))
	for i, member := range members {
		interfaceMembers[i] = member
	}

	return s.rdb.SRem(ctx, key, interfaceMembers...).Err()
}

func (s *setOperator) SCard(ctx context.Context, key string) (int64, error) {
	result, err := s.rdb.SCard(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get set cardinality")
	}

	return result, nil
}
