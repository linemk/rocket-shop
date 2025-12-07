package http

import (
	"context"
	"net/http"
	"strings"
)

const (
	SessionUUIDHeader = "X-Session-UUID"
)

type contextKey string

const sessionUUIDContextKey contextKey = "session-uuid"

// AuthMiddleware возвращает middleware для проверки сессии в HTTP запросах
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			http.Error(w, "Missing session uuid header", http.StatusUnauthorized)
			return
		}

		// Добавляем session UUID в контекст
		ctx := context.WithValue(r.Context(), sessionUUIDContextKey, sessionUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware возвращает middleware, который опционально проверяет сессию
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.Header.Get(SessionUUIDHeader)

		ctx := r.Context()
		if sessionUUID != "" {
			ctx = context.WithValue(ctx, sessionUUIDContextKey, sessionUUID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ExtractSessionUUID извлекает session UUID из контекста
func ExtractSessionUUID(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDContextKey).(string)
	return sessionUUID, ok
}

// AddSessionUUIDToRequest добавляет session UUID в заголовок запроса
func AddSessionUUIDToRequest(r *http.Request, sessionUUID string) {
	r.Header.Set(SessionUUIDHeader, sessionUUID)
}

// ForwardSessionUUIDToGRPC добавляет session UUID в gRPC metadata string
func ForwardSessionUUIDToGRPC(ctx context.Context) string {
	if sessionUUID, ok := ExtractSessionUUID(ctx); ok {
		return sessionUUID
	}
	return ""
}

// ExtractSessionUUIDFromHeader извлекает session UUID из заголовка строки
func ExtractSessionUUIDFromHeader(headerValue string) string {
	// Формат: "Bearer <session_uuid>"
	parts := strings.Fields(headerValue)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return headerValue
}
