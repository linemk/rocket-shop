package api

import (
	"context"

	authv1 "github.com/linemk/rocket-shop/shared/pkg/proto/auth/v1"

	"github.com/linemk/rocket-shop/iam/internal/service/auth"
	"github.com/linemk/rocket-shop/iam/internal/service/converter"
)

type authV1Handler struct {
	authService auth.Service
	authv1.UnimplementedAuthServiceServer
}

func NewAuthV1Handler(authService auth.Service) authv1.AuthServiceServer {
	return &authV1Handler{
		authService: authService,
	}
}

func (h *authV1Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	session, err := h.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &authv1.LoginResponse{
		SessionUuid: session.SessionUUID,
	}, nil
}

func (h *authV1Handler) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	user, err := h.authService.Whoami(ctx, req.SessionUuid)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &authv1.WhoamiResponse{
		User: converter.UserToProto(user),
	}, nil
}

func (h *authV1Handler) handleError(err error) error {
	// TODO: Map domain errors to gRPC status codes
	return err
}
