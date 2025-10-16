package apperrors

import "errors"

var (
	ErrPartNotFound      = errors.New("part not found")
	ErrPartOutOfStock    = errors.New("part is out of stock")
	ErrInvalidFilter     = errors.New("invalid filter")
	ErrInvalidUUID       = errors.New("invalid UUID format")
	ErrPartAlreadyExists = errors.New("part already exists")
)
