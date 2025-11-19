package decoder

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/linemk/rocket-shop/notification/internal/model"
	events_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/events/v1"
)

// DecodeShipAssembled декодирует событие ShipAssembled из JSON
func DecodeShipAssembled(data []byte) (*model.ShipAssembledEvent, error) {
	var protoEvent events_v1.ShipAssembled

	// Пробуем protojson
	if err := protojson.Unmarshal(data, &protoEvent); err != nil {
		// Если не получилось, пробуем обычный JSON
		if err := json.Unmarshal(data, &protoEvent); err != nil {
			return nil, err
		}
	}

	return &model.ShipAssembledEvent{
		EventUUID:    protoEvent.EventUuid,
		OrderUUID:    protoEvent.OrderUuid,
		UserUUID:     protoEvent.UserUuid,
		BuildTimeSec: protoEvent.BuildTimeSec,
	}, nil
}
