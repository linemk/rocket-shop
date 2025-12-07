package model

import "time"

// Session представляет сессию пользователя
type Session struct {
	SessionUUID string
	UserUUID    string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}
