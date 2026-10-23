package order

import (
	"context"
	"errors"

	"github.com/ivmerk/goshop/order/internal/model"
	repoConverter "github.com/ivmerk/goshop/order/internal/repository/converter"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

func (r *repository) Create(ctx context.Context, order *model.Order) (string, error) {

	if order.UUID == "" {
		return "", model.ErrInvalidOrderUUID
	}

	o := repoConverter.ToRepositoryOrder(order)

	_, err := r.pool.Exec(ctx,
		`INSERT INTO orders (uuid, "user", parts, total, payment, transaction, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		o.UUID, o.User, o.Parts, o.Total, o.Payment, o.Transaction, o.Status)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return "", model.ErrOrderAlreadyExists
		}
		return "", err
	}
	return o.UUID, nil
}
