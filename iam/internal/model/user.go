package model

// NotificationMethod представляет метод уведомления пользователя
type NotificationMethod struct {
	ProviderName string
	Target       string
}

// User представляет пользователя в системе
type User struct {
	UserUUID            string
	Login               string
	PasswordHash        string
	Email               string
	NotificationMethods []NotificationMethod
}
