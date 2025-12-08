package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/linemk/rocket-shop/iam/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, userUUID string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}
