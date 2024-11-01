package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "angkot",
	Short: "Telegram bot application",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	viper.AutomaticEnv()
}
