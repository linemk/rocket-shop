package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/linemk/rocket-shop/inventory/internal/converter"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func (a *API) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	partInfo, err := a.inventoryUseCase.GetPart(ctx, req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Part with UUID %s not found", req.Uuid)
	}

	protoPart := converter.PartInfoToProto(partInfo)
	return &inventory_v1.GetPartResponse{
		Part: protoPart,
	}, nil
}
