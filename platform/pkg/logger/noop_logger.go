package logger

import (
	"context"

	"go.uber.org/zap"
)

// NoopLogger - логгер без операций для использования в тестах и closer
type NoopLogger struct{}

func (n *NoopLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {}
