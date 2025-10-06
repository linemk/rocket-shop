package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/linemk/rocket-shop"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderService struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

type Order struct {
	ID            string   `json:"order_id"`
	UserID        string   `json:"user_id"`
	DetailsID     []string `json:"details_id"`
	TotalPrice    float64  `json:"total_price"`
	TransactionID string   `json:"transaction_id"`
	PaymentMethod string   `json:"payment_method"`
	Status        string   `json:"status"`
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*Order),
	}
}

func (s *OrderService) GetOrder(id string) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, ok := s.orders[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("order by id %s not found", id))
	}

	return w, nil
}

func GetOrder(ctx context.Context) (GetOrderRes, error) {

}

func main() {

}
