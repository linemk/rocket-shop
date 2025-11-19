package app

import (
	"context"
	"log"

	"github.com/linemk/rocket-shop/notification/internal/config"
	"github.com/linemk/rocket-shop/platform/pkg/closer"
)

type App struct {
	diContainer *diContainer
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		if err := closer.CloseAll(ctx); err != nil {
			log.Printf("failed to close all resources: %s", err.Error())
		}
		closer.Wait()
	}()

	// Запускаем оба Kafka consumers параллельно
	return a.runConsumers(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initDiContainer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %s", err.Error())
	}

	return nil
}

func (a *App) initDiContainer(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) runConsumers(ctx context.Context) error {
	// Запускаем оба consumer'а одновременно
	errChan := make(chan error, 2)

	go func() {
		errChan <- a.diContainer.OrderPaidConsumer(ctx).RunConsumer(ctx)
	}()

	go func() {
		errChan <- a.diContainer.OrderAssembledConsumer(ctx).RunConsumer(ctx)
	}()

	// Ждем первой ошибки
	return <-errChan
}
