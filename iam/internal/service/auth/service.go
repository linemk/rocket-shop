package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/linemk/rocket-shop/iam/internal/config"
	"github.com/linemk/rocket-shop/iam/internal/model"
	"github.com/linemk/rocket-shop/iam/internal/repository/session"
	"github.com/linemk/rocket-shop/iam/internal/repository/user"
)

type Service interface {
	Login(ctx context.Context, login, password string) (*model.Session, error)
	Whoami(ctx context.Context, sessionUUID string) (*model.User, error)
}

type service struct {
	userRepo    user.Repository
	sessionRepo session.Repository
	sessionCfg  config.SessionConfig
}

func NewService(userRepo user.Repository, sessionRepo session.Repository, sessionCfg config.SessionConfig) Service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionCfg:  sessionCfg,
	}
}

func (s *service) Login(ctx context.Context, login, password string) (*model.Session, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, errors.Wrap(err, "failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(s.sessionCfg.TTL())

	session := &model.Session{
		SessionUUID: uuid.New().String(),
		UserUUID:    user.UserUUID,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
	}

	err = s.sessionRepo.Create(ctx, session, s.sessionCfg.TTL())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	err = s.sessionRepo.AddSessionToUserSet(ctx, user.UserUUID, session.SessionUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add session to user set")
	}

	return session, nil
}

func (s *service) Whoami(ctx context.Context, sessionUUID string) (*model.User, error) {
	session, err := s.sessionRepo.Get(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, model.ErrSessionNotFound
		}
		return nil, errors.Wrap(err, "failed to get session")
	}

	if time.Now().After(session.ExpiresAt) {
		err := s.sessionRepo.Delete(ctx, sessionUUID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete expired session")
		}
		return nil, model.ErrSessionExpired
	}

	user, err := s.userRepo.GetByID(ctx, session.UserUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return user, nil
}
