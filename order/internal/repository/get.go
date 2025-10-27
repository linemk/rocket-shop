package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (models.Order, error) {
	query := `SELECT
		uuid,
		user_id,
		part_uuids,
		total_price,
		transaction_id,
		payment_method,
		status,
		created_at,
		updated_at
	FROM orders
	WHERE uuid = $1`

	var order models.Order
	var partUUIDs []uuid.UUID
	var paymentMethodStr, statusStr string
	var updatedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, orderUUID).Scan(
		&order.UUID,
		&order.UserID,
		&partUUIDs, // pgx автоматически конвертирует PostgreSQL array в []uuid.UUID
		&order.TotalPrice,
		&order.TransactionID,
		&paymentMethodStr,
		&statusStr,
		&order.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, apperrors.ErrOrderNotFound
		}
		return models.Order{}, err
	}

	// Присваиваем прочитанные значения
	order.PartUUIDs = partUUIDs
	order.PaymentMethod = order_v1.PaymentMethod(paymentMethodStr)
	order.Status = order_v1.OrderStatus(statusStr)

	if updatedAt.Valid {
		order.UpdatedAt = &updatedAt.Time
	}

	return order, nil
}
