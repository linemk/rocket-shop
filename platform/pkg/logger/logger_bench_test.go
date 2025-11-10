package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

// BenchmarkLogger_Info бенчмарк для измерения производительности логгера
func BenchmarkLogger_Info(b *testing.B) {
	InitForBenchmark()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(ctx, "benchmark message", zap.String("key", "value"))
	}
}

// BenchmarkLogger_InfoWithContext бенчмарк с контекстом
func BenchmarkLogger_InfoWithContext(b *testing.B) {
	InitForBenchmark()
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(ctx, "benchmark message", zap.String("key", "value"))
	}
}
