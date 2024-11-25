package commands

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robzlabz/angkot-kaf/database"
	"github.com/robzlabz/angkot-kaf/models"
)

func HandleReport(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var targetDate time.Time
	var err error

	args := message.CommandArguments()
	if args == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("02-01-2006", args)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Format tanggal tidak valid. Gunakan format DD-MM-YYYY")
			bot.Send(msg)
			return
		}
	}

	// Get all trips for the target date
	var trips []models.Trip
	err = database.DB.
		Preload("Driver").
		Preload("Passengers").
		Where("DATE(trip_date) = DATE(?)", targetDate).
		Find(&trips).Error

	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Gagal mengambil data perjalanan.")
		bot.Send(msg)
		return
	}

	if len(trips) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("📊 Tidak ada perjalanan pada tanggal %s", targetDate.Format("02-01-2006")))
		bot.Send(msg)
		return
	}

	// Group trips by driver and passenger
	driverTrips := make(map[uint][]models.Trip)
	passengerTrips := make(map[string][]models.Trip)

	for _, trip := range trips {
		driverTrips[trip.DriverID] = append(driverTrips[trip.DriverID], trip)
		for _, passenger := range trip.Passengers {
			passengerTrips[passenger.Name] = append(passengerTrips[passenger.Name], trip)
		}
	}

	// Calculate total revenue
	var totalRevenue int
	response := fmt.Sprintf("📊 Laporan Perjalanan %s\n\n", targetDate.Format("Monday, 02 January 2006"))

	// Process each driver's trips
	for _, dTrips := range driverTrips {
		if len(dTrips) == 0 {
			continue
		}

		response += fmt.Sprintf("🚗 Driver: %s\n", dTrips[0].Driver.Name)

		// Process antar trips
		antarTrips := filterTripsByType(dTrips, "antar")
		if len(antarTrips) > 0 {
			response += "\n🔜 Antar:\n"
			for _, trip := range antarTrips {
				for _, passenger := range trip.Passengers {
					response += fmt.Sprintf("- %s\n", passenger.Name)
				}
			}
		}

		// Process jemput trips
		jemputTrips := filterTripsByType(dTrips, "jemput")
		if len(jemputTrips) > 0 {
			response += "\n🔙 Jemput:\n"
			for _, trip := range jemputTrips {
				for _, passenger := range trip.Passengers {
					response += fmt.Sprintf("- %s\n", passenger.Name)
				}
			}
		}

		response += "\n"
	}

	// Add payment summary
	response += "💰 Ringkasan Pembayaran:\n\n"

	for passenger, pTrips := range passengerTrips {
		var hasAntar, hasJemput bool
		for _, t := range pTrips {
			if t.TripType == "antar" {
				hasAntar = true
			} else if t.TripType == "jemput" {
				hasJemput = true
			}
		}

		fare := OneWayFare
		tripType := "satu arah"
		if hasAntar && hasJemput {
			fare = TwoWayFare
			tripType = "pulang pergi"
		}

		totalRevenue += fare
		response += fmt.Sprintf("👤 %s (%s)", passenger, tripType)
		response += fmt.Sprintf(" 💵 Rp %d\n", fare)
	}

	response += fmt.Sprintf("\n💵 Total Pendapatan: Rp %d", totalRevenue)

	// Send report in chunks if it's too long
	if len(response) > 4096 {
		chunks := splitMessage(response)
		for _, chunk := range chunks {
			msg := tgbotapi.NewMessage(message.Chat.ID, chunk)
			bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
	}
}

func filterTripsByType(trips []models.Trip, tripType string) []models.Trip {
	var filtered []models.Trip
	for _, trip := range trips {
		if trip.TripType == tripType {
			filtered = append(filtered, trip)
		}
	}
	return filtered
}

func splitMessage(message string) []string {
	var chunks []string
	for len(message) > 0 {
		if len(message) <= 4096 {
			chunks = append(chunks, message)
			break
		}

		// Find the last newline within the 4096 character limit
		chunk := message[:4096]
		lastNewline := strings.LastIndex(chunk, "\n")
		if lastNewline == -1 {
			lastNewline = 4096
		}

		chunks = append(chunks, message[:lastNewline])
		message = message[lastNewline+1:]
	}
	return chunks
}
