package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/linemk/rocket-shop/iam/internal/model"
	repoConverter "github.com/linemk/rocket-shop/iam/internal/repository/converter"
	repoModel "github.com/linemk/rocket-shop/iam/internal/repository/model"
	"github.com/linemk/rocket-shop/platform/pkg/cache/redis"
)

func (r *repository) Get(ctx context.Context, sessionUUID string) (*model.Session, error) {
	key := fmt.Sprintf("%s%s", sessionKeyPrefix, sessionUUID)

	sessionJSON, err := r.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.ErrKeyNotFound) {
			return nil, model.ErrSessionNotFound
		}
		return nil, errors.Wrap(err, "failed to get session from Redis")
	}

	var repoSession repoModel.Session
	err = json.Unmarshal(sessionJSON, &repoSession)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal session")
	}

	return repoConverter.ToInternalSession(&repoSession), nil
}
