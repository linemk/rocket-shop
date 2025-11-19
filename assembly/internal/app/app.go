package app

import (
	"context"
	"log"

	"github.com/linemk/rocket-shop/assembly/internal/config"
	"github.com/linemk/rocket-shop/platform/pkg/closer"
	"github.com/linemk/rocket-shop/platform/pkg/logger"
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

	// Запускаем Kafka consumers
	return a.runConsumers(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initCloser,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(config.AppConfig().Logger.Level(), false)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initDiContainer(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) runConsumers(ctx context.Context) error {
	return a.diContainer.ConsumerService(ctx).RunConsumers(ctx)
}
