//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ivmerk/goshop/platform/pkg/logger"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

type InventorySuite struct {
	suite.Suite

	env    *TestEnvironment
	conn   *grpc.ClientConn
	client inventoryV1.InventoryServiceClient
}

func TestInventoryE2E(t *testing.T) {
	suite.Run(t, new(InventorySuite))
}

func (s *InventorySuite) SetupSuite() {
	// platform/testcontainers пишет через глобальный логгер при ошибках
	s.Require().NoError(logger.Init("info", false))

	env, err := setupTestEnvironment(context.Background())
	s.Require().NoError(err)
	s.env = env

	conn, err := grpc.NewClient(env.App.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = inventoryV1.NewInventoryServiceClient(conn)
}

func (s *InventorySuite) TearDownSuite() {
	if s.conn != nil {
		s.NoError(s.conn.Close())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if s.env != nil {
		teardownTestEnvironment(ctx, s.env)
	}
}

var _ = require.NoError // убрать, если require не нужен в этом файле
