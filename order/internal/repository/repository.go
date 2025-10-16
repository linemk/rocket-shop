package repository

import (
	"context"
	"sync"

	"github.com/linemk/rocket-shop/order/internal/entyties/models"
)

type repository struct {
	mu   sync.RWMutex
	data map[string]models.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]models.Order),
	}
}

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) error
	Get(ctx context.Context, uuid string) (models.Order, error)
	Update(ctx context.Context, uuid string, updateInfo models.OrderUpdateInfo) error
}
