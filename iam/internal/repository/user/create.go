package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/linemk/rocket-shop/iam/internal/model"
	repoConverter "github.com/linemk/rocket-shop/iam/internal/repository/converter"
	repoModel "github.com/linemk/rocket-shop/iam/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, user *model.User) error {
	repoUser := repoConverter.ToRepoUser(user)
	if repoUser == nil {
		return errors.New("failed to convert user to repository model")
	}

	notificationMethodsJSON, err := repoModel.NotificationMethodsToJSON(repoUser.NotificationMethods)
	if err != nil {
		return errors.Wrap(err, "failed to marshal notification methods")
	}

	now := time.Now()

	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(
			"user_uuid",
			"login",
			"password_hash",
			"email",
			"notification_methods",
			"created_at",
		).
		Values(
			repoUser.UserUUID,
			repoUser.Login,
			repoUser.PasswordHash,
			repoUser.Email,
			notificationMethodsJSON,
			now,
		).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	return nil
}
