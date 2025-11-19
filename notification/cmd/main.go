package main

import (
	"context"
	"log"

	"github.com/linemk/rocket-shop/notification/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	log.Println("NotificationService is starting...")

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
