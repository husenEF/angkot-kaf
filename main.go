package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/robzlabz/angkot/internal/infrastructure/bot"
	"github.com/robzlabz/angkot/pkg/logging"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Starting bot...")
	logging.InitLogging()

	// Load .env file
	fmt.Println("Loading .env file...")
	if err := godotenv.Load(); err != nil {
		logging.Error("Error loading .env file", err)
	}

	fmt.Println("Loading environment variables...")
	viper.AutomaticEnv()

	fmt.Println("Starting bot...")
	err := bot.Start()
	if err != nil {
		logging.Error("[Main] Bot Error", err)
	}

}
