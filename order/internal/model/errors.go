package model

import "errors"

var (
	ErrOrderAlreadyExists   = errors.New("order already exists")
	ErrOrderNotFound        = errors.New("order not found")
	ErrInvalidOrderUUID     = errors.New("invalid order uuid")
	ErrEmptyParts           = errors.New("empty parts list")
	ErrPartNotFound         = errors.New("parts not found")
	ErrInvalidPaymentMethod = errors.New("wrong payment method")
	ErrOrderAlreadyPaid     = errors.New("order already paid")
	ErrOrderNotPayable      = errors.New("")
)
