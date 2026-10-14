package model

import "errors"

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrPaymentMethodUnknown = errors.New("unknown payment method")
)
