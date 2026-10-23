package order

import (
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"
)

const (
	pendingOrderUUID = "11111111-1111-1111-1111-111111111111"
	paidOrderUUID    = "22222222-2222-2222-2222-222222222222"
	userUUID         = "33333333-3333-3333-3333-333333333333"
	partUUID1        = "44444444-4444-4444-4444-444444444444"
	partUUID2        = "55555555-5555-5555-5555-555555555555"
	transactionUUID  = "66666666-6666-6666-6666-666666666666"
)

func newTestRepo(t *testing.T) (*repository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	t.Cleanup(mock.Close)
	return &repository{pool: mock}, mock
}
