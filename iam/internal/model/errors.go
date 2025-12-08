package model

import "errors"

var (
	// ErrUserNotFound возвращается когда пользователь не найден
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists возвращается когда пользователь с таким логином уже существует
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials возвращается при неверных учетных данных
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrSessionNotFound возвращается когда сессия не найдена
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired возвращается когда сессия истекла
	ErrSessionExpired = errors.New("session expired")

	// ErrInvalidSessionUUID возвращается когда session UUID невалиден
	ErrInvalidSessionUUID = errors.New("invalid session UUID")

	// ErrInvalidUserUUID возвращается когда user UUID невалиден
	ErrInvalidUserUUID = errors.New("invalid user UUID")
)
