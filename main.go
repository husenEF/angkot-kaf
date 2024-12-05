package main

import (
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/robzlabz/angkot-kaf/commands"
	"github.com/robzlabz/angkot-kaf/database"
	"github.com/robzlabz/angkot-kaf/models"
)

func init() {
	// Set timezone to Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Warning: Failed to load Asia/Jakarta timezone: %v. Falling back to UTC+7", err)
		// Create a fixed timezone offset for Jakarta (UTC+7)
		loc = time.FixedZone("WIB", 7*60*60) // UTC+7
	}
	time.Local = loc

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	database.Init()

	// Auto-migrate database models
	if err := database.DB.AutoMigrate(&models.Driver{}, &models.Trip{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set bot commands
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Mulai bot dan lihat menu utama"},
		{Command: "add_driver", Description: "Mendaftarkan driver baru"},
		{Command: "gas", Description: "Mencatat perjalanan antar/jemput"},
		{Command: "report", Description: "Melihat catatan hari ini atau tanggal tertentu"},
		{Command: "backupdb", Description: "Backup database"},
	}

	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Printf("Error setting bot commands: %v", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)
	handleUpdates(bot, updates)
}

func handleUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		go handleMessage(bot, update.Message)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleMessage: %v", r)
			msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Terjadi kesalahan internal. Silakan coba lagi.")
			bot.Send(msg)
		}
	}()

	if message.IsCommand() {
		handleCommand(bot, message)
		return
	}

	// Handle trip messages (antar/jemput)
	messageText := strings.ToLower(message.Text)
	if strings.HasPrefix(messageText, "antar") || strings.HasPrefix(messageText, "jemput") {
		commands.HandleTrip(bot, message)
	}

}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		commands.HandleStart(bot, message)
	case "add_driver":
		commands.HandleAddDriver(bot, message)
	case "gas":
		commands.HandleGas(bot, message)
	case "report":
		commands.HandleReport(bot, message)
	case "backupdb":
		commands.HandleBackupDB(bot, message)
	case "delete":
		commands.HandleDeleteTrip(bot, message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Perintah tidak dikenali. Gunakan /start untuk melihat daftar perintah.")
		bot.Send(msg)
	}
}
