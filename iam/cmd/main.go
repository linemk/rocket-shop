package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/linemk/rocket-shop/iam/internal/app"
	"github.com/linemk/rocket-shop/platform/pkg/closer"
	"github.com/linemk/rocket-shop/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app: %v\n", err)
		os.Exit(1)
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan

		logger.Info(ctx, fmt.Sprintf("Received signal: %v, shutting down...", sig))

		if err := closer.CloseAll(ctx); err != nil {
			logger.Error(ctx, fmt.Sprintf("Error during shutdown: %v", err))
		}
	}()

	if err := application.Run(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("Application error: %v", err))
		os.Exit(1)
	}
}
