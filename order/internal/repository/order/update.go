package order

import (
	"context"

	"github.com/ivmerk/goshop/order/internal/model"
)

func (r *repository) Update(ctx context.Context, uuid string, updatedOrder *model.OrderUpdateInfo) (string, error) {

	tag, err := r.pool.Exec(ctx,
		`UPDATE orders 
	SET  transaction = COALESCE($2, transaction),
		status      = COALESCE(NULLIF($3, ''), status),
    payment     = COALESCE(NULLIF($4, ''), payment),
    updated_at  = now()
	WHERE uuid = $1`, uuid, updatedOrder.Transaction, string(updatedOrder.Status), string(updatedOrder.Payment))
	if err != nil {
		return "", err
	}
	if tag.RowsAffected() == 0 {
		return "", model.ErrOrderNotFound
	}
	return uuid, nil

}
