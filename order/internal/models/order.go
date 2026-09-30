package models

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSpb           PaymentMethod = "SPB"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет собой структуру данных для хранения информации о заказе.
type Order struct {
	ID          string        `json:"id"`          // Уникальный идентификатор заказа
	User        string        `json:"user"`        // Имя клиента, сделавшего заказ
	Parts       []string      `json:"parts"`       // Список товаров в заказе
	Total       float64       `json:"total"`       // Общая стоимость заказа
	Transaction string        `json:"transaction"` // Идентификатор транзакции оплаты заказа
	Status      OrderStatus   `json:"status"`      // Статус заказа (например, "PENDING_PAYMENT")
	Payment     PaymentMethod `json:"payment"`     // Метод оплаты заказа (например, "CARD", "SPB", "CREDIT_CARD")
	CreatedAt   time.Time     `json:"createdAt"`   // Время создания заказа
}
