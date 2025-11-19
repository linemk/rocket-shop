package model

// OrderPaidEvent представляет событие об оплате заказа
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

// ShipAssembledEvent представляет событие о завершении сборки корабля
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
