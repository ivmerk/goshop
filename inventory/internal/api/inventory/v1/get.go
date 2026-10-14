package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivmerk/goshop/inventory/internal/converter"
	"github.com/ivmerk/goshop/inventory/internal/model"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)}, nil
}
