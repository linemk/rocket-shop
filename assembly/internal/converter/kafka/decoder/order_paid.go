package decoder

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/linemk/rocket-shop/assembly/internal/model"
	events_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/events/v1"
)

// DecodeOrderPaid декодирует событие OrderPaid из JSON
func DecodeOrderPaid(data []byte) (*model.OrderPaidEvent, error) {
	var protoEvent events_v1.OrderPaid

	// Пробуем сначала protojson
	if err := protojson.Unmarshal(data, &protoEvent); err != nil {
		// Если не получилось, пробуем обычный JSON
		if err := json.Unmarshal(data, &protoEvent); err != nil {
			return nil, err
		}
	}

	return &model.OrderPaidEvent{
		EventUUID:       protoEvent.EventUuid,
		OrderUUID:       protoEvent.OrderUuid,
		UserUUID:        protoEvent.UserUuid,
		PaymentMethod:   protoEvent.PaymentMethod,
		TransactionUUID: protoEvent.TransactionUuid,
	}, nil
}
