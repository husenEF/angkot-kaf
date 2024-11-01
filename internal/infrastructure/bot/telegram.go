package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robzlabz/angkot/internal/core/services"
	"github.com/robzlabz/angkot/internal/infrastructure/database"
	"github.com/spf13/viper"
)

func Start() {
	sqlDB, err := database.NewSQLiteDB()
	if err != nil {
		log.Printf("[Adapter][Start]Error initializing SQLite database: %v", err)
		log.Fatal(err)
	}

	// Use SQLite as primary storage
	botService := services.NewBotService(sqlDB)

	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_TOKEN"))
	if err != nil {
		log.Printf("[Adapter][Start]Error initializing Telegram bot: %v", err)
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("[Adapter][Start]Error getting updates channel: %v", err)
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messageText := update.Message.Text

		switch {
		case messageText == "ping":
			response := botService.HandlePing()
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "penumpang":
			response := botService.HandlePassenger(chatID)
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "daftarpenumpang":
			response, err := botService.GetPassengerList()
			if err != nil {
				response = "Maaf, terjadi kesalahan saat membaca data penumpang"
			}
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "driver":
			response := botService.HandleDriver(chatID)
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "drivers":
			response, err := botService.GetDriverList()
			if err != nil {
				response = "Maaf, terjadi kesalahan saat membaca data driver"
			}
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		default:
			if strings.HasPrefix(strings.ToLower(messageText), "keberangkatan") {
				response, err := botService.ProcessDeparture(messageText)
				if err != nil {
					response = "Maaf, terjadi kesalahan saat memproses keberangkatan"
				}
				msg := tgbotapi.NewMessage(chatID, response)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
				}
			} else if strings.HasPrefix(strings.ToLower(messageText), "kepulangan") {
				response, err := botService.ProcessReturn(messageText)
				if err != nil {
					response = "Maaf, terjadi kesalahan saat memproses kepulangan"
				}
				msg := tgbotapi.NewMessage(chatID, response)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
				}
			} else if botService.IsWaitingForPassengerName(chatID) {
				err := botService.AddPassenger(messageText)
				var response string
				if err != nil {
					response = "Maaf, terjadi kesalahan saat menyimpan data penumpang"
				} else {
					response = "Penumpang " + messageText + " berhasil ditambahkan"
				}
				botService.ClearWaitingStatus(chatID)
				msg := tgbotapi.NewMessage(chatID, response)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
				}
			} else if botService.IsWaitingForDriverName(chatID) {
				err := botService.AddDriver(messageText)
				var response string
				if err != nil {
					response = "Maaf, terjadi kesalahan saat menyimpan data driver"
				} else {
					response = "Driver " + messageText + " berhasil ditambahkan"
				}
				botService.ClearWaitingStatus(chatID)
				msg := tgbotapi.NewMessage(chatID, response)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
				}
			}
		}
	}
}
