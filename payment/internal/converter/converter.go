package converter

import (
	"github.com/linemk/rocket-shop/payment/internal/entyties/models"
)

// TransactionToProto конвертирует модель Transaction в protobuf (если бы он существовал)
// Пока что просто возвращаем nil, так как в protobuf нет типа Transaction
func TransactionToProto(transaction models.Transaction) interface{} {
	// В protobuf схеме нет типа Transaction, поэтому возвращаем nil
	return nil
}

// ProtoToTransaction конвертирует protobuf Transaction в модель Transaction
// Пока что просто возвращаем пустую структуру, так как в protobuf нет типа Transaction
func ProtoToTransaction(protoTransaction interface{}) models.Transaction {
	// В protobuf схеме нет типа Transaction, поэтому возвращаем пустую структуру
	return models.Transaction{}
}
