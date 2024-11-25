package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, 
		"🚐 Selamat datang di Angkot KAF Bot!\n\n" +
		"Gunakan perintah berikut:\n" +
		"📝 /add_driver - Mendaftarkan driver baru\n" +
		"🚗 /gas - Mencatat perjalanan\n" +
		"📊 /report - Melihat catatan hari ini\n" +
		"📊 /report DD-MM-YYYY - Melihat catatan tanggal tertentu\n" +
		"💾 /backupdb - Backup database")
	bot.Send(msg)
}
