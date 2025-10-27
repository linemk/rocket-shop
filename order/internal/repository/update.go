package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

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

	// Используем squirrel для type-safe query building
	query, args, err := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", order.UpdatedAt).
		Set("status", string(order.Status)).
		Set("transaction_id", order.TransactionID).
		Set("payment_method", string(order.PaymentMethod)).
		Where(sq.Eq{"uuid": uuid}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrOrderNotFound
	}

	return nil
}
