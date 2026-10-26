package order

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	def "github.com/ivmerk/goshop/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type dbPool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type repository struct {
	pool dbPool
}

func NewOrderRepository(pool dbPool) *repository {
	return &repository{
		pool: pool,
	}
}
