//go:build integration

package integration

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

func (s *InventorySuite) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (s *InventorySuite) TestHealthCheck() {
	ctx, cancel := s.ctx()
	defer cancel()

	resp, err := grpc_health_v1.NewHealthClient(s.conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	s.Require().NoError(err)
	s.Equal(grpc_health_v1.HealthCheckResponse_SERVING, resp.GetStatus())
}

func (s *InventorySuite) TestListParts_ReturnsSeededParts() {
	ctx, cancel := s.ctx()
	defer cancel()

	resp, err := s.client.ListParts(ctx, &inventoryV1.ListPartsRequest{})
	s.Require().NoError(err)
	s.Len(resp.GetParts(), 4) // Seed создаёт 4 детали
}

func (s *InventorySuite) TestGetPart_Found() {
	ctx, cancel := s.ctx()
	defer cancel()

	list, err := s.client.ListParts(ctx, &inventoryV1.ListPartsRequest{})
	s.Require().NoError(err)
	s.Require().NotEmpty(list.GetParts())
	want := list.GetParts()[0]

	got, err := s.client.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: want.GetUuid()})
	s.Require().NoError(err)
	s.Equal(want.GetUuid(), got.GetPart().GetUuid())
	s.Equal(want.GetName(), got.GetPart().GetName())
}

func (s *InventorySuite) TestGetPart_NotFound() {
	ctx, cancel := s.ctx()
	defer cancel()

	_, err := s.client.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: "00000000-0000-0000-0000-000000000000"})
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *InventorySuite) TestListParts_FilterByUUID() {
	ctx, cancel := s.ctx()
	defer cancel()

	all, err := s.client.ListParts(ctx, &inventoryV1.ListPartsRequest{})
	s.Require().NoError(err)
	target := all.GetParts()[0].GetUuid()

	resp, err := s.client.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: []string{target}},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.GetParts(), 1)
	s.Equal(target, resp.GetParts()[0].GetUuid())
}
