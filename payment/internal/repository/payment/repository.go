package payment

import (
	"context"
	"fmt"
	"sync"

	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
)

type Repository struct {
	mu           sync.RWMutex
	transactions map[string]models.Transaction
}

func NewRepository() *Repository {
	return &Repository{
		transactions: make(map[string]models.Transaction),
	}
}

func (r *Repository) CreateTransaction(ctx context.Context, transaction models.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.UUID]; exists {
		return fmt.Errorf("transaction with UUID %s already exists", transaction.UUID)
	}

	r.transactions[transaction.UUID] = transaction
	return nil
}

func (r *Repository) GetTransaction(ctx context.Context, uuid string) (models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transaction, exists := r.transactions[uuid]
	if !exists {
		return models.Transaction{}, fmt.Errorf("transaction with UUID %s not found", uuid)
	}

	return transaction, nil
}

func (r *Repository) UpdateTransaction(ctx context.Context, uuid string, transaction models.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[uuid]; !exists {
		return fmt.Errorf("transaction with UUID %s not found", uuid)
	}

	r.transactions[uuid] = transaction
	return nil
}

func (r *Repository) ListTransactions(ctx context.Context, orderUUID string) ([]models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.Transaction
	for _, transaction := range r.transactions {
		if transaction.OrderUUID == orderUUID {
			result = append(result, transaction)
		}
	}

	return result, nil
}
