package v1

import (
	"context"
	"net/http"

	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) *order_v1.UnexpectedErrStatusCode {
	return &order_v1.UnexpectedErrStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.UnexpectedErr{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
