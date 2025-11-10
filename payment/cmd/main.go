package main

import (
	"context"
	"log"
	"syscall"

	"github.com/linemk/rocket-shop/payment/internal/app"
	"github.com/linemk/rocket-shop/platform/pkg/closer"
)

func main() {
	ctx := context.Background()

	// Создаем приложение
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Настраиваем graceful shutdown
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	// Запускаем приложение в горутине
	go func() {
		if err := application.Run(ctx); err != nil {
			log.Fatalf("Failed to run application: %v", err)
		}
	}()

	// Ожидаем завершения
	closer.Wait()
}
