package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robzlabz/angkot/internal/core/services"
	"github.com/robzlabz/angkot/internal/infrastructure/database"
	"github.com/spf13/viper"
)

func Start() error {
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
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("[Adapter][Start]Error getting updates channel: %v", err)
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messageText := update.Message.Text

		fmt.Printf("Received message: %s from chat ID: %d\n", messageText, chatID)

		switch {
		case messageText == "/ping":
			response := botService.HandlePing()
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/santri":
			response := botService.HandlePassenger(chatID)
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/daftarsantri":
			response, err := botService.GetPassengerList()
			if err != nil {
				response = "Maaf, terjadi kesalahan saat membaca data santri"
			}
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/driver":
			response := botService.HandleDriver(chatID)
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/drivers":
			response, err := botService.GetDriverList()
			if err != nil {
				response = "Maaf, terjadi kesalahan saat membaca data driver"
			}
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/keberangkatan":
			response := "Format pencatatan keberangkatan:\n\n" +
				"keberangkatan\n" +
				"Driver: [nama_driver]\n" +
				"- [nama_santri_1]\n" +
				"- [nama_santri_2]\n" +
				"- [nama_santri_3]\n\n" +
				"Contoh:\n" +
				"keberangkatan\n" +
				"Driver: Pak Ahmad\n" +
				"- Santri Ali\n" +
				"- Santri Umar\n" +
				"- Santri Hasan"
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/kepulangan":
			response := "Format pencatatan kepulangan:\n\n" +
				"kepulangan\n" +
				"Driver: [nama_driver]\n" +
				"- [nama_santri_1]\n" +
				"- [nama_santri_2]\n" +
				"- [nama_santri_3]\n\n" +
				"Contoh:\n" +
				"kepulangan\n" +
				"Driver: Pak Ahmad\n" +
				"- Santri Ali\n" +
				"- Santri Umar\n" +
				"- Santri Hasan"
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/start" || messageText == "/help":
			response := "Selamat datang! Berikut adalah daftar perintah yang tersedia:\n" +
				"/ping - Cek koneksi bot\n" +
				"/driver - Tambah driver baru\n" +
				"/drivers - Lihat daftar driver\n" +
				"/antar - Lihat format pencatatan antar\n" +
				"/jemput - Lihat format pencatatan jemput\n" +
				"/laporan - Lihat laporan harian"
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/laporan":
			response, err := botService.GetTodayReport()
			if err != nil {
				response = "Maaf, terjadi kesalahan saat membuat laporan"
			}
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/antar":
			response := "Format pencatatan antar:\n\n" +
				"antar\n" +
				"Driver: [nama_driver]\n" +
				"- [nama_santri_1]\n" +
				"- [nama_santri_2]\n" +
				"- [nama_santri_3]\n\n" +
				"Contoh:\n" +
				"antar\n" +
				"Driver: Pak Ahmad\n" +
				"- Santri Ali\n" +
				"- Santri Umar\n" +
				"- Santri Hasan"
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		case messageText == "/jemput":
			response := "Format pencatatan jemput:\n\n" +
				"jemput\n" +
				"Driver: [nama_driver]\n" +
				"- [nama_santri_1]\n" +
				"- [nama_santri_2]\n" +
				"- [nama_santri_3]\n\n" +
				"Contoh:\n" +
				"jemput\n" +
				"Driver: Pak Ahmad\n" +
				"- Santri Ali\n" +
				"- Santri Umar\n" +
				"- Santri Hasan"
			msg := tgbotapi.NewMessage(chatID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
			}
		default:
			if strings.HasPrefix(strings.ToLower(messageText), "antar") {
				lines := strings.Split(messageText, "\n")
				response, err := botService.ProcessDeparture(lines[1], lines[2:])
				if err != nil {
					response = err.Error()
				}
				msg := tgbotapi.NewMessage(chatID, response)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("[Adapter][MessageHandler]Error sending message: %v", err)
				}
			} else if strings.HasPrefix(strings.ToLower(messageText), "jemput") {
				lines := strings.Split(messageText, "\n")
				response, err := botService.ProcessReturn(lines[1], lines[2:])
				if err != nil {
					response = err.Error()
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

	return nil
}
