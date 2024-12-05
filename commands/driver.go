package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robzlabz/angkot-kaf/database"
	"github.com/robzlabz/angkot-kaf/models"
)

func HandleAddDriver(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "ℹ️ Gunakan format: /add_driver <nama_driver>")
		bot.Send(msg)
		return
	}

	driver := models.Driver{
		Name:   args,
		Active: true,
	}

	if err := database.DB.Create(&driver).Error; err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Gagal menambahkan driver. Silakan coba lagi.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "✅ Driver berhasil ditambahkan!")
	bot.Send(msg)
}

func HandleGas(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID,
		"📝 Format pencatatan perjalanan:\n\n"+
			"1️⃣ Format Antar:\n"+
			"antar\n"+
			"driver: <nama_driver>\n"+
			"date: <tanggal> (opsional,eg: 25-03-2024)\n"+
			"1. <nama_penumpang_1>\n"+
			"2. <nama_penumpang_2>\n\n"+
			"2️⃣ Format Jemput:\n"+
			"jemput\n"+
			"driver: <nama_driver>\n"+
			"date: <tanggal> (opsional,eg: 25-03-2024)\n"+
			"1. <nama_penumpang_1>\n"+
			"2. <nama_penumpang_2>\n\n"+
			"3️⃣ Format Report Tanggal Tertentu:\n"+
			"report DD-MM-YYYY\n"+
			"Contoh: report 25-03-2024\n\n"+
			"💰 Biaya:\n"+
			"- Antar saja: Rp 10.000\n"+
			"- Antar + Jemput: Rp 15.000")
	bot.Send(msg)
}
