package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const SessionUUIDHeader = "session-uuid"

type contextKey string

const sessionUUIDContextKey contextKey = "session-uuid"

// UnaryAuthInterceptor возвращает unary interceptor для проверки сессии
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not found")
	}

	sessionUUIDs := md.Get(SessionUUIDHeader)
	if len(sessionUUIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "session uuid not found")
	}

	// Добавляем session UUID в контекст
	ctx = context.WithValue(ctx, sessionUUIDContextKey, sessionUUIDs[0])

	return handler(ctx, req)
}

// StreamAuthInterceptor возвращает stream interceptor для проверки сессии
func StreamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata not found")
	}

	sessionUUIDs := md.Get(SessionUUIDHeader)
	if len(sessionUUIDs) == 0 {
		return status.Error(codes.Unauthenticated, "session uuid not found")
	}

	// Создаем новый контекст с session UUID
	newCtx := context.WithValue(ss.Context(), sessionUUIDContextKey, sessionUUIDs[0])
	return handler(srv, &contextWrappedStream{stream: ss, ctx: newCtx})
}

type contextWrappedStream struct {
	stream grpc.ServerStream
	ctx    context.Context //nolint:containedctx
}

func (w *contextWrappedStream) SetHeader(md metadata.MD) error {
	return w.stream.SetHeader(md)
}

func (w *contextWrappedStream) SendHeader(md metadata.MD) error {
	return w.stream.SendHeader(md)
}

func (w *contextWrappedStream) SetTrailer(md metadata.MD) {
	w.stream.SetTrailer(md)
}

func (w *contextWrappedStream) Context() context.Context {
	return w.ctx
}

func (w *contextWrappedStream) SendMsg(m interface{}) error {
	return w.stream.SendMsg(m)
}

func (w *contextWrappedStream) RecvMsg(m interface{}) error {
	return w.stream.RecvMsg(m)
}

// ExtractSessionUUID извлекает session UUID из контекста
func ExtractSessionUUID(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDContextKey).(string)
	return sessionUUID, ok
}
