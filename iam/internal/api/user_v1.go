package api

import (
	"context"

	userv1 "github.com/linemk/rocket-shop/shared/pkg/proto/user/v1"

	"github.com/linemk/rocket-shop/iam/internal/model"
	"github.com/linemk/rocket-shop/iam/internal/service/converter"
	userservice "github.com/linemk/rocket-shop/iam/internal/service/user"
)

type userV1Handler struct {
	userService userservice.Service
	userv1.UnimplementedUserServiceServer
}

func NewUserV1Handler(userService userservice.Service) userv1.UserServiceServer {
	return &userV1Handler{
		userService: userService,
	}
}

func (h *userV1Handler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	notificationMethods := make([]model.NotificationMethod, len(req.NotificationMethods))
	for i, method := range req.NotificationMethods {
		notificationMethods[i] = model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	user, err := h.userService.Register(ctx, req.Login, req.Password, req.Email, notificationMethods)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &userv1.RegisterResponse{
		UserUuid: user.UserUUID,
	}, nil
}

func (h *userV1Handler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := h.userService.GetUser(ctx, req.UserUuid)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &userv1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}

func (h *userV1Handler) handleError(err error) error {
	// TODO: Map domain errors to gRPC status codes
	return err
}
