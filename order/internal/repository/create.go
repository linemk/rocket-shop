package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

func (r *repository) Create(ctx context.Context, order models.Order) error {
	now := time.Now()

	// Используем squirrel для type-safe query building
	query, args, err := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns(
			"uuid",
			"user_id",
			"part_uuids",
			"total_price",
			"transaction_id",
			"payment_method",
			"status",
			"created_at",
		).
		Values(
			order.UUID,
			order.UserID,
			order.PartUUIDs, // pgx автоматически конвертирует []uuid.UUID в PostgreSQL array
			order.TotalPrice,
			order.TransactionID,
			string(order.PaymentMethod),
			string(order.Status),
			now,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}
