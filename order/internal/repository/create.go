package repository

import (
	"context"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
	"time"
)

func (r *repository) Create(_ context.Context, order models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = nil

	r.data[order.UUID] = order
	return nil
}
