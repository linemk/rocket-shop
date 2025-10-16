package apperrors

import "errors"

var (
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrInvalidPaymentMethod     = errors.New("invalid payment method")
	ErrPaymentFailed            = errors.New("payment failed")
	ErrTransactionAlreadyExists = errors.New("transaction already exists")
	ErrInvalidAmount            = errors.New("invalid amount")
)
