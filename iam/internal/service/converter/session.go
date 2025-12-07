package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "github.com/linemk/rocket-shop/shared/pkg/proto/common/v1"

	"github.com/linemk/rocket-shop/iam/internal/model"
)

func SessionToProto(session *model.Session) *commonv1.Session {
	if session == nil {
		return nil
	}

	return &commonv1.Session{
		SessionUuid: session.SessionUUID,
		UserUuid:    session.UserUUID,
		CreatedAt:   timestamppb.New(session.CreatedAt),
		ExpiresAt:   timestamppb.New(session.ExpiresAt),
	}
}

func SessionFromProto(session *commonv1.Session) *model.Session {
	if session == nil {
		return nil
	}

	return &model.Session{
		SessionUUID: session.SessionUuid,
		UserUUID:    session.UserUuid,
		CreatedAt:   session.CreatedAt.AsTime(),
		ExpiresAt:   session.ExpiresAt.AsTime(),
	}
}
