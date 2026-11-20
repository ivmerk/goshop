// Package metrics — бизнес-метрики Order-сервиса.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// OrdersTotal — количество созданных заказов (на нём построен алерт «>10 заказов за минуту»).
	OrdersTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_total",
		Help: "Общее количество созданных заказов.",
	})

	// OrdersRevenueTotal — выручка: сумма стоимости оплаченных заказов.
	OrdersRevenueTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_revenue_total",
		Help: "Суммарная выручка по оплаченным заказам.",
	})
)
