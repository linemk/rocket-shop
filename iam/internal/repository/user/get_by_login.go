package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/linemk/rocket-shop/iam/internal/model"
	repoConverter "github.com/linemk/rocket-shop/iam/internal/repository/converter"
	repoModel "github.com/linemk/rocket-shop/iam/internal/repository/model"
)

func (r *repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return r.queryUser(ctx, sq.Eq{"login": login})
}

func (r *repository) queryUser(ctx context.Context, where sq.Eq) (*model.User, error) {
	query, args, err := sq.Select(
		"user_uuid",
		"login",
		"password_hash",
		"email",
		"notification_methods",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(where).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}

	var repoUser repoModel.User
	var notificationMethodsJSON []byte

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&repoUser.UserUUID,
		&repoUser.Login,
		&repoUser.PasswordHash,
		&repoUser.Email,
		&notificationMethodsJSON,
		&repoUser.CreatedAt,
		&repoUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "failed to scan user")
	}

	repoUser.NotificationMethods, err = repoModel.NotificationMethodsFromJSON(notificationMethodsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal notification methods")
	}

	return repoConverter.ToInternalUser(&repoUser), nil
}
