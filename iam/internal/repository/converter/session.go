package converter

import (
	internalModel "github.com/linemk/rocket-shop/iam/internal/model"
	repoModel "github.com/linemk/rocket-shop/iam/internal/repository/model"
)

// ToInternalSession конвертирует repository Session в internal Session
func ToInternalSession(session *repoModel.Session) *internalModel.Session {
	if session == nil {
		return nil
	}

	return &internalModel.Session{
		SessionUUID: session.SessionUUID,
		UserUUID:    session.UserUUID,
		CreatedAt:   session.CreatedAt,
		ExpiresAt:   session.ExpiresAt,
	}
}

// ToRepoSession конвертирует internal Session в repository Session
func ToRepoSession(session *internalModel.Session) *repoModel.Session {
	if session == nil {
		return nil
	}

	return &repoModel.Session{
		SessionUUID: session.SessionUUID,
		UserUUID:    session.UserUUID,
		CreatedAt:   session.CreatedAt,
		ExpiresAt:   session.ExpiresAt,
	}
}
