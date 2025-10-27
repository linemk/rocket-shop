package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) error
	Get(ctx context.Context, uuid string) (models.Order, error)
	Update(ctx context.Context, uuid string, updateInfo models.OrderUpdateInfo) error
}
