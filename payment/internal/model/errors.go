package model

import "errors"

var (
	ErrInvalidOrderUUID         = errors.New("invalid order UUID")
	ErrInvalidUserUUID          = errors.New("invalid user UUID")
	ErrUnsupportedPaymentMethod = errors.New("unsupported payment method")
)
