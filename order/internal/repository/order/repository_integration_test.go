//go:build integration

package order

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/ivmerk/goshop/order/internal/model"
	"github.com/ivmerk/goshop/platform/pkg/migrator"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/path"
)

const (
	pgImage    = "postgres:17.0-alpine3.20" // та же версия, что в docker-compose
	pgUser     = "order_user"
	pgPassword = "order_password"
	pgDatabase = "order"
	pgPort     = "5432/tcp"

	migrationsRelPath = "order/migrations"

	orderUUID1  = "11111111-1111-1111-1111-111111111111"
	orderUUID2  = "22222222-2222-2222-2222-222222222222"
	missingUUID = "99999999-9999-9999-9999-999999999999"
	userName    = "user-1"

// partUUID1       = "44444444-4444-4444-4444-444444444444"
// partUUID2       = "55555555-5555-5555-5555-555555555555"
// transactionUUID = "66666666-6666-6666-6666-666666666666"
)

type OrderRepoSuite struct {
	suite.Suite

	container testcontainers.Container
	pool      *pgxpool.Pool
	repo      *repository
}

func TestOrderRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(OrderRepoSuite))
}

func (s *OrderRepoSuite) SetupSuite() {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pgImage,
			ExposedPorts: []string{pgPort},
			Env: map[string]string{
				"POSTGRES_USER":     pgUser,
				"POSTGRES_PASSWORD": pgPassword,
				"POSTGRES_DB":       pgDatabase,
			},
			// Postgres при инициализации стартует дважды: ждём второго «ready»,
			// иначе wait.ForListeningPort срабатывает на временный сервер.
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Minute),
		},
		Started: true,
	})
	s.Require().NoError(err)
	s.container = container

	host, err := container.Host(ctx)
	s.Require().NoError(err)
	port, err := container.MappedPort(ctx, pgPort)
	s.Require().NoError(err)

	uri := "postgres://" + pgUser + ":" + pgPassword + "@" + host + ":" + port.Port() +
		"/" + pgDatabase + "?sslmode=disable"

	s.pool, err = pgxpool.New(ctx, uri)
	s.Require().NoError(err)
	s.Require().NoError(s.pool.Ping(ctx))

	// настоящие миграции, не копия схемы
	migrationsDir := path.GetProjectRoot() + "/" + migrationsRelPath
	s.Require().NoError(migrator.NewMigrator(stdlib.OpenDBFromPool(s.pool), migrationsDir).Up())

	s.repo = NewOrderRepository(s.pool)
}

func (s *OrderRepoSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.container != nil {
		s.NoError(s.container.Terminate(context.Background()))
	}
}

// Каждый тест начинает с пустой таблицы.
func (s *OrderRepoSuite) SetupTest() {
	_, err := s.pool.Exec(context.Background(), "TRUNCATE TABLE orders")
	s.Require().NoError(err)
}

func newOrder(uuid string) *model.Order {
	return &model.Order{
		UUID:    uuid,
		User:    userUUID,
		Parts:   []string{partUUID1, partUUID2},
		Total:   3200.5,
		Status:  model.OrderStatusPendingPayment,
		Payment: model.PaymentMethodCard,
	}
}

func (s *OrderRepoSuite) TestCreateAndGet() {
	ctx := context.Background()
	want := newOrder(orderUUID1)

	id, err := s.repo.Create(ctx, want)
	s.Require().NoError(err)
	s.Equal(orderUUID1, id)

	got, err := s.repo.Get(ctx, orderUUID1)
	s.Require().NoError(err)
	s.Equal(want.UUID, got.UUID)
	s.Equal(want.User, got.User)
	s.Equal(want.Parts, got.Parts)
	s.InDelta(want.Total, got.Total, 0.001)
	s.Equal(want.Status, got.Status)
	s.Equal(want.Payment, got.Payment)
	s.Nil(got.Transaction)
}

func (s *OrderRepoSuite) TestCreate_Duplicate() {
	ctx := context.Background()

	_, err := s.repo.Create(ctx, newOrder(orderUUID1))
	s.Require().NoError(err)

	_, err = s.repo.Create(ctx, newOrder(orderUUID1))
	s.ErrorIs(err, model.ErrOrderAlreadyExists)
}

func (s *OrderRepoSuite) TestGet_NotFound() {
	_, err := s.repo.Get(context.Background(), missingUUID)
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *OrderRepoSuite) TestUpdate_StatusAndTransaction() {
	ctx := context.Background()

	_, err := s.repo.Create(ctx, newOrder(orderUUID1))
	s.Require().NoError(err)

	tx := transactionUUID
	_, err = s.repo.Update(ctx, orderUUID1, &model.OrderUpdateInfo{
		Status:      model.OrderStatusPaid,
		Payment:     model.PaymentMethodSbp,
		Transaction: &tx,
	})
	s.Require().NoError(err)

	got, err := s.repo.Get(ctx, orderUUID1)
	s.Require().NoError(err)
	s.Equal(model.OrderStatusPaid, got.Status)
	s.Equal(model.PaymentMethodSbp, got.Payment)
	s.Require().NotNil(got.Transaction)
	s.Equal(transactionUUID, *got.Transaction)

	var updatedAt *time.Time
	s.Require().NoError(s.pool.QueryRow(ctx,
		"SELECT updated_at FROM orders WHERE uuid = $1", orderUUID1).Scan(&updatedAt))
	s.NotNil(updatedAt, "updated_at должен проставляться при Update")
}

func (s *OrderRepoSuite) TestUpdate_EmptyFieldsKeepOldValues() {
	ctx := context.Background()

	_, err := s.repo.Create(ctx, newOrder(orderUUID1))
	s.Require().NoError(err)

	_, err = s.repo.Update(ctx, orderUUID1, &model.OrderUpdateInfo{}) // всё пустое
	s.Require().NoError(err)

	got, err := s.repo.Get(ctx, orderUUID1)
	s.Require().NoError(err)
	s.Equal(model.OrderStatusPendingPayment, got.Status) // COALESCE/NULLIF не затёрли значения
	s.Equal(model.PaymentMethodCard, got.Payment)
}

func (s *OrderRepoSuite) TestUpdate_NotFound() {
	_, err := s.repo.Update(context.Background(), missingUUID, &model.OrderUpdateInfo{Status: model.OrderStatusPaid})
	s.ErrorIs(err, model.ErrOrderNotFound)
}
