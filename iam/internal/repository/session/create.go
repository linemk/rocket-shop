package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/linemk/rocket-shop/iam/internal/model"
	repoConverter "github.com/linemk/rocket-shop/iam/internal/repository/converter"
)

const sessionKeyPrefix = "session:"

func (r *repository) Create(ctx context.Context, session *model.Session, ttl time.Duration) error {
	repoSession := repoConverter.ToRepoSession(session)
	if repoSession == nil {
		return errors.New("failed to convert session to repository model")
	}

	sessionJSON, err := json.Marshal(repoSession)
	if err != nil {
		return errors.Wrap(err, "failed to marshal session")
	}

	key := fmt.Sprintf("%s%s", sessionKeyPrefix, repoSession.SessionUUID)

	err = r.cache.Set(ctx, key, sessionJSON, ttl)
	if err != nil {
		return errors.Wrap(err, "failed to save session to Redis")
	}

	return nil
}
