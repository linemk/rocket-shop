package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

type PartInfo struct {
	UUID          string
	Name          string
	Price         float32
	StockQuantity int64
}

func (c *Client) GetPart(ctx context.Context, partUUID uuid.UUID) (PartInfo, error) {
	resp, err := c.client.GetPart(ctx, &inventory_v1.GetPartRequest{
		Uuid: partUUID.String(),
	})
	if err != nil {
		return PartInfo{}, fmt.Errorf("failed to get part %s: %w", partUUID.String(), err)
	}

	return PartInfo{
		UUID:          resp.Part.Uuid,
		Name:          resp.Part.Name,
		Price:         float32(resp.Part.Price),
		StockQuantity: resp.Part.StockQuantity,
	}, nil
}
