package model

import (
	"database/sql"
	"encoding/json"
)

// NotificationMethod представляет метод уведомления в БД
type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

// User представляет пользователя в БД
type User struct {
	UserUUID            string
	Login               string
	PasswordHash        string
	Email               string
	NotificationMethods []NotificationMethod
	CreatedAt           sql.NullTime
	UpdatedAt           sql.NullTime
}

// NotificationMethodsToJSON конвертирует NotificationMethods в JSONB для PostgreSQL
func NotificationMethodsToJSON(methods []NotificationMethod) ([]byte, error) {
	if len(methods) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(methods)
}

// NotificationMethodsFromJSON парсит JSONB из PostgreSQL в NotificationMethods
func NotificationMethodsFromJSON(data []byte) ([]NotificationMethod, error) {
	var methods []NotificationMethod
	if len(data) == 0 {
		return methods, nil
	}
	err := json.Unmarshal(data, &methods)
	return methods, err
}
