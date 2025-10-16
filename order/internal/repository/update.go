package repository

import (
	"context"
	"time"

	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

func (r *repository) Update(_ context.Context, uuid string, updateInfo models.OrderUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.data[uuid]
	if !exists {
		return apperrors.ErrOrderNotFound
	}

	now := time.Now()
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

	r.data[uuid] = order
	return nil
}
