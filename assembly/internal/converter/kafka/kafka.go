package kafka

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/linemk/rocket-shop/assembly/internal/model"
	events_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/events/v1"
)

// EncodeShipAssembled кодирует событие ShipAssembled в JSON
func EncodeShipAssembled(event *model.ShipAssembledEvent) ([]byte, error) {
	protoEvent := &events_v1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	// Используем protojson для маршалинга
	data, err := protojson.Marshal(protoEvent)
	if err != nil {
		// Fallback на обычный JSON
		return json.Marshal(protoEvent)
	}

	return data, nil
}
