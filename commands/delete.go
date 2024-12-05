package commands

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robzlabz/angkot-kaf/database"
	"github.com/robzlabz/angkot-kaf/models"
)

func HandleDeleteTrip(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Parse the command (e.g., "/delete antar 2024-12-05" or "/delete antar")
	parts := strings.Fields(message.Text)
	if len(parts) < 2 {
		sendMessage(bot, message.Chat.ID, "❌ Format tidak valid. Gunakan: /delete [antar/jemput] [YYYY-MM-DD (opsional)]")
		return
	}

	tripType := strings.ToLower(parts[1])
	if tripType != "antar" && tripType != "jemput" {
		sendMessage(bot, message.Chat.ID, "❌ Tipe perjalanan tidak valid. Gunakan 'antar' atau 'jemput'.")
		return
	}

	// If a date is provided (i.e., len(parts) == 3), parse it. Otherwise, use current date.
	var tripDate time.Time
	var err error
	if len(parts) == 3 {
		// Custom date provided by user
		tripDate, err = time.Parse("02/01/2006", parts[2])
		if err != nil {
			sendMessage(bot, message.Chat.ID, "❌ Format tanggal tidak valid. Gunakan format: DD/MM/YYYY.")
			return
		}
	} else {
		// No date provided, use current date
		tripDate = time.Now()
	}

	// Delete trips for the specified date and type
	if err := database.DB.Where("trip_type = ? AND DATE(trip_date) = ?", tripType, tripDate.Format("2006-01-02")).
		Delete(&models.Trip{}).Error; err != nil {
		sendMessage(bot, message.Chat.ID, "❌ Gagal menghapus data perjalanan.")
		return
	}

	// Send confirmation message
	sendMessage(bot, message.Chat.ID, fmt.Sprintf("✅ Semua perjalanan dengan tipe '%s' untuk tanggal %s telah dihapus.", tripType, tripDate.Format("02-01-2006")))
}
