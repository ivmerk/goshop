package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
	repoModel "github.com/ivmerk/goshop/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Order, error) {
	var o repoModel.Order
	err := r.pool.QueryRow(ctx,
		`SELECT uuid, "user", parts, total, payment, transaction, status, created_at, updated_at FROM orders WHERE uuid=$1`, uuid).Scan(&o.UUID, &o.User, &o.Parts, &o.Total, &o.Payment,
		&o.Transaction, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}
	return repoConverter.FromRepositoryOrder(&o), nil
}
