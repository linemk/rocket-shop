package model

import "time"

// Session представляет сессию в Redis
type Session struct {
	SessionUUID string
	UserUUID    string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}
