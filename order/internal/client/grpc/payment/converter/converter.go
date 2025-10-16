package converter

import (
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

func OpenAPIPaymentMethodToProto(openAPIMethod order_v1.PaymentMethod) payment_v1.PaymentMethod {
	switch openAPIMethod {
	case order_v1.PaymentMethodPAYMENTMETHODUNSPECIFIED:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	case order_v1.PaymentMethodPAYMENTMETHODCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case order_v1.PaymentMethodPAYMENTMETHODSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case order_v1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case order_v1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
