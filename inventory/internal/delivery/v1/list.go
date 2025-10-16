package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/linemk/rocket-shop/inventory/internal/converter"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

func (a *API) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	filter := converter.ProtoToPartFilter(req.Filter)

	partInfos, err := a.inventoryUseCase.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list parts: %v", err)
	}

	// Конвертируем []PartInfo в []*inventory_v1.Part
	protoParts := make([]*inventory_v1.Part, 0, len(partInfos))
	for _, partInfo := range partInfos {
		protoPart := converter.PartInfoToProto(partInfo)
		protoParts = append(protoParts, protoPart)
	}

	return &inventory_v1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
