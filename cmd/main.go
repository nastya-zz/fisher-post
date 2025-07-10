package main

import (
	"context"
	"post/internal/app"
	"post/pkg/logger"
)

func main() {
	logger.Init()

	ctx := context.Background()

	logger.Info("Starting auth service")
	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to init app", "error", err.Error())
	}
	if a == nil {
		logger.Fatal("failed to init app")
		return
	}

	err = a.Run(ctx)
	if err != nil {
		logger.Fatal("failed to run app", "error", err.Error())
	}
}
