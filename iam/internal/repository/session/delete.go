package session

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (r *repository) Delete(ctx context.Context, sessionUUID string) error {
	key := fmt.Sprintf("%s%s", sessionKeyPrefix, sessionUUID)

	err := r.cache.Del(ctx, key)
	if err != nil {
		return errors.Wrap(err, "failed to delete session from Redis")
	}

	return nil
}
