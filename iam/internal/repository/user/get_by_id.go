package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/linemk/rocket-shop/iam/internal/model"
)

func (r *repository) GetByID(ctx context.Context, userUUID string) (*model.User, error) {
	return r.queryUser(ctx, sq.Eq{"user_uuid": userUUID})
}
