package repository

import (
	"context"

	"github.com/linemk/rocket-shop/order/internal/entyties/apperrors"
	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

func (r *repository) Get(_ context.Context, uuid string) (models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.data[uuid]
	if !exists {
		return models.Order{}, apperrors.ErrOrderNotFound
	}

	return order, nil
}
