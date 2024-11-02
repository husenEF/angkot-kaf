package main

import (
	"github.com/joho/godotenv"
	"github.com/robzlabz/angkot/internal/infrastructure/bot"
	"github.com/robzlabz/angkot/pkg/logging"
	"github.com/spf13/viper"
)

func main() {
	logging.InitLogging()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logging.Error("Error loading .env file", err)
	}

	viper.AutomaticEnv()

	err := bot.Start()
	logging.Info("[Main] Bot started")
	if err != nil {
		logging.Error("[Main] Bot Error", err)
	}

}
