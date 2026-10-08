package model

import "time"

type OrderStatus string
type PaymentMethod string

type Order struct {
	UUID        string        `json:"uuid"`        // Уникальный идентификатор заказа
	User        string        `json:"user"`        // Имя клиента, сделавшего заказ
	Parts       []string      `json:"parts"`       // Список товаров в заказе
	Total       float64       `json:"total"`       // Общая стоимость заказа
	Payment     PaymentMethod `json:"payment"`     // Метод оплаты заказа (например, "CARD", "SBP", "CREDIT_CARD")
	Transaction *string       `json:"transaction"` // Идентификатор транзакции оплаты заказа
	Status      OrderStatus   `json:"status"`      // Статус заказа (например, "PENDING_PAYMENT")
	CreatedAt   time.Time     `json:"created_at"`  // Время создания заказа
	UpdatedAt   *time.Time    `json:"updated_at"`  // Время последнего обновления заказа
}
