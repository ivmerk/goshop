package v1

import (
	service "github.com/ivmerk/goshop/inventory/internal/service"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewAPI(partService service.PartService) *api {
	return &api{partService: partService}
}
