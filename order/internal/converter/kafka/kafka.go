package kafka

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/linemk/rocket-shop/order/internal/entyties/events"
	events_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/events/v1"
)

// EncodeOrderPaid кодирует событие OrderPaid в JSON
func EncodeOrderPaid(event *events.OrderPaidEvent) ([]byte, error) {
	protoEvent := &events_v1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	// Используем protojson для маршалинга
	data, err := protojson.Marshal(protoEvent)
	if err != nil {
		// Fallback на обычный JSON
		return json.Marshal(protoEvent)
	}

	return data, nil
}
