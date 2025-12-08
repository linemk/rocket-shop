package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/linemk/rocket-shop/iam/internal/model"
	"github.com/linemk/rocket-shop/iam/internal/repository/user"
)

type Service interface {
	Register(ctx context.Context, login, password, email string, notificationMethods []model.NotificationMethod) (*model.User, error)
	GetUser(ctx context.Context, userUUID string) (*model.User, error)
}

type service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) Register(ctx context.Context, login, password, email string, notificationMethods []model.NotificationMethod) (*model.User, error) {
	existingUser, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		return nil, errors.Wrap(err, "failed to check if user exists")
	}

	if existingUser != nil {
		return nil, model.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	newUser := &model.User{
		UserUUID:            uuid.New().String(),
		Login:               login,
		PasswordHash:        string(passwordHash),
		Email:               email,
		NotificationMethods: notificationMethods,
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return newUser, nil
}

func (s *service) GetUser(ctx context.Context, userUUID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return user, nil
}
