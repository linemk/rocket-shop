package repository

import (
	"context"
	"time"

	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

func (r *repository) Create(ctx context.Context, order models.Order) error {
	now := time.Now()

	// Используем прямой INSERT с pgx для работы с массивами
	query := `INSERT INTO orders
	(uuid, user_id, part_uuids, total_price, transaction_id, payment_method, status, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(ctx, query,
		order.UUID,
		order.UserID,
		order.PartUUIDs, // pgx автоматически конвертирует []uuid.UUID в PostgreSQL array
		order.TotalPrice,
		order.TransactionID,
		string(order.PaymentMethod),
		string(order.Status),
		now,
	)

	return err
}
