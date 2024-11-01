package cmd

import (
	"github.com/robzlabz/angkot/internal/infrastructure/bot"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the telegram bot",
	Run: func(cmd *cobra.Command, args []string) {
		bot.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
