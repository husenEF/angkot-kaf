package commands

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleBackupDB(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	now := time.Now()
	backupFileName := fmt.Sprintf("backup_angkot_%s.db", now.Format("2006-01-02_150405"))

	// Read the original database file
	dbData, err := os.ReadFile("database/angkot.db")
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Gagal membaca database.")
		bot.Send(msg)
		return
	}

	// Write to backup file
	err = os.WriteFile(backupFileName, dbData, 0644)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Gagal membuat file backup.")
		bot.Send(msg)
		return
	}

	// Send the backup file
	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FilePath(backupFileName))
	_, err = bot.Send(doc)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Gagal mengirim file backup.")
		bot.Send(msg)
		return
	}

	// Clean up the backup file
	os.Remove(backupFileName)
}
