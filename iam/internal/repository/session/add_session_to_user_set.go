package session

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

const userSessionsKeyPrefix = "user_sessions:"

func (r *repository) AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	key := fmt.Sprintf("%s%s", userSessionsKeyPrefix, userUUID)

	err := r.cache.SetOperator().SAdd(ctx, key, sessionUUID)
	if err != nil {
		return errors.Wrap(err, "failed to add session to user set")
	}

	return nil
}
