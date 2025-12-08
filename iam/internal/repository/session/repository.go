package session

import (
	"context"
	"time"

	"github.com/linemk/rocket-shop/iam/internal/model"
	"github.com/linemk/rocket-shop/platform/pkg/cache"
)

type Repository interface {
	Create(ctx context.Context, session *model.Session, ttl time.Duration) error
	Get(ctx context.Context, sessionUUID string) (*model.Session, error)
	Delete(ctx context.Context, sessionUUID string) error
	AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error
}

type repository struct {
	cache cache.Client
}

func NewRepository(cache cache.Client) Repository {
	return &repository{
		cache: cache,
	}
}
