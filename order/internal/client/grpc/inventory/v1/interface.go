package v1

import (
	"context"

	"github.com/google/uuid"
)

type InventoryClient interface {
	GetPart(ctx context.Context, partUUID uuid.UUID) (PartInfo, error)
	Close() error
}
