package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivmerk/goshop/inventory/internal/converter"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &inventoryV1.ListPartsResponse{Parts: converter.PartsToProto(parts)}, nil
}
