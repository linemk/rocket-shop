package apperrors

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidParts      = errors.New("invalid parts")
	ErrPartOutOfStock    = errors.New("part out of stock")
	ErrPaymentFailed     = errors.New("payment failed")
	ErrOrderCannotBePaid = errors.New("order cannot be paid in current status")
	ErrOrderAlreadyPaid  = errors.New("order already paid and cannot be cancelled")
	ErrNoPartsSpecified  = errors.New("no parts specified")
	ErrPartNotFound      = errors.New("part not found in inventory")
)
