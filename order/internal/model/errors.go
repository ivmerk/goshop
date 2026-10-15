package model

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderUUID   = errors.New("invalid order uuid")
)
