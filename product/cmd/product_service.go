package main

import (
	"context"

	"github.com/Puena/auction/logger"
	"github.com/Puena/auction/product/internal/app"
)

func main() {
	// create app config
	appConfig, err := app.NewConfig(app.WithEnv())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed while creating app config")
	}

	// init app with config
	app := app.NewApp(appConfig)

	// run app database migration
	err = app.RunMigration()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed while running migration")
	}

	// run app
	err = app.Run(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed while running app")
	}
}
