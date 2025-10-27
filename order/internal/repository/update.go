package repository

import (
	"context"
	"time"

	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

func (r *repository) Update(ctx context.Context, uuid string, updateInfo models.OrderUpdateInfo) error {
	now := time.Now()

	// Получаем текущий заказ
	order, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}

	// Обновляем поля
	order.UpdatedAt = &now
	if updateInfo.Status != nil {
		order.Status = *updateInfo.Status
	}
	if updateInfo.TransactionID != nil {
		order.TransactionID = *updateInfo.TransactionID
	}
	if updateInfo.PaymentMethod != nil {
		order.PaymentMethod = *updateInfo.PaymentMethod
	}

	// Обновляем в БД
	query := `UPDATE orders
	SET updated_at = $1, status = $2, transaction_id = $3, payment_method = $4
	WHERE uuid = $5`

	result, err := r.db.Exec(ctx, query,
		order.UpdatedAt,
		string(order.Status),
		order.TransactionID,
		string(order.PaymentMethod),
		uuid,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrOrderNotFound
	}

	return nil
}
