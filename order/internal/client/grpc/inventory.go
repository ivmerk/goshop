package grpc

import (
	"context"
	"fmt"

	inventory_v1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

type InventoryClient struct {
	client inventory_v1.InventoryServiceClient
}

func NewInventoryClient(client inventory_v1.InventoryServiceClient) *InventoryClient {
	return &InventoryClient{client: client}
}

// ListPartPrices возвращает цены найденных деталей: ключ — UUID детали.
// Не найденные детали в результат не попадают.
func (c *InventoryClient) ListPartPrices(ctx context.Context, uuids []string) (map[string]float64, error) {
	resp, err := c.client.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{Uuids: uuids},
	})
	if err != nil {
		return nil, fmt.Errorf("inventory ListParts: %w", err)
	}

	prices := make(map[string]float64, len(resp.GetParts()))
	for _, part := range resp.GetParts() {
		prices[part.GetUuid()] = part.GetPrice()
	}

	return prices, nil
}
