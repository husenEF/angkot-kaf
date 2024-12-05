package commands

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robzlabz/angkot-kaf/database"
	"github.com/robzlabz/angkot-kaf/models"
)

const (
	OneWayFare = 10000 // Rp 10.000 for one-way trip
	TwoWayFare = 15000 // Rp 15.000 for round trip
)

type tripInfo struct {
	driverName string
	passengers []string
	tripType   string
	tripDate   time.Time
}

func HandleTrip(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	info, err := parseTripMessage(message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("❌ %s", err.Error()))
		bot.Send(msg)
		return
	}

	// Find driver
	var driver models.Driver
	if err := database.DB.Where("name = ? AND active = ?", info.driverName, true).First(&driver).Error; err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Driver tidak ditemukan atau tidak aktif.")
		bot.Send(msg)
		return
	}

	// Save trips and passengers
	var response string
	for _, passengerName := range info.passengers {
		passenger := models.Passenger{
			Name:     passengerName,
			TripType: info.tripType,
			Amount:   OneWayFare, // Default to one-way fare
		}

		// Create trip record with passenger
		trip := models.Trip{
			DriverID:   driver.ID,
			Driver:     driver,
			TripType:   info.tripType,
			TripDate:   info.tripDate,
			Passengers: []models.Passenger{passenger},
		}

		if err := database.DB.Create(&trip).Error; err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("❌ Gagal mencatat perjalanan untuk %s", passengerName))
			bot.Send(msg)
			continue
		}
	}

	if info.tripType == "antar" {
		response = formatAntarResponse(driver.Name, info.passengers, info.tripDate)
	} else {
		response = formatJemputResponse(driver.Name, info.passengers, time.Now())
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func parseTripMessage(text string) (*tripInfo, error) {
	lines := strings.Split(text, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("Format tidak valid. Gunakan /gas untuk melihat format yang benar.")
	}

	tripType := strings.ToLower(strings.TrimSpace(lines[0]))
	if tripType != "antar" && tripType != "jemput" {
		return nil, fmt.Errorf("Tipe perjalanan tidak valid. Gunakan 'antar' atau 'jemput'.")
	}

	// Parse driver
	driverLine := strings.TrimSpace(lines[1])
	if !strings.HasPrefix(strings.ToLower(driverLine), "driver:") {
		return nil, fmt.Errorf("Format driver tidak valid. Gunakan 'driver: nama_driver'.")
	}
	driverName := strings.TrimSpace(strings.TrimPrefix(driverLine, "driver:"))

	// Check for date in the second line
	var tripDate time.Time
	var passengerLine int = 2
	// Parse date
	dateLine := strings.TrimSpace(lines[2])
	if strings.HasPrefix(strings.ToLower(dateLine), "date:") {
		passengerLine = 3
		dateStr := strings.TrimSpace(strings.TrimPrefix(dateLine, "date:"))
		dateLayout := "02/01/2006"
		datetime, err := time.Parse(dateLayout, dateStr)
		if err != nil {
			return nil, fmt.Errorf("Format tanggal tidak valid. Gunakan format DD/MM/YYYY")
		}
		tripDate = datetime
	} else {
		tripDate = time.Now()
	}

	// Parse passengers
	var passengers []string
	for i := passengerLine; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ".", 2)
		if len(parts) != 2 {
			continue
		}

		passengerName := strings.TrimSpace(parts[1])
		if passengerName != "" {
			passengers = append(passengers, passengerName)
		}
	}

	if len(passengers) == 0 {
		return nil, fmt.Errorf("Tidak ada penumpang yang valid.")
	}

	return &tripInfo{
		driverName: driverName,
		passengers: passengers,
		tripType:   tripType,
		tripDate:   tripDate,
	}, nil
}

func formatAntarResponse(driverName string, passengers []string, tripTime time.Time) string {
	response := fmt.Sprintf("✅ Berhasil mencatat perjalanan antar:\nDriver: %s", driverName)
	response += fmt.Sprintf("\nDate: %s\n", tripTime.Format("02-01-2006"))
	response += "\nPenumpang:\n"
	for _, passenger := range passengers {
		response += fmt.Sprintf("- %s\n", passenger)
	}
	return response
}

func formatJemputResponse(driverName string, passengers []string, tripTime time.Time) string {
	response := fmt.Sprintf("✅ Berhasil mencatat perjalanan jemput:\nDriver: %s", driverName)
	response += fmt.Sprintf("\nDate: %s\n", tripTime.Format("02-01-2006"))
	response += "\nRingkasan Pembayaran:\n"

	// Calculate fares for each passenger
	for _, passenger := range passengers {
		// Check if passenger has both antar and jemput trips today
		var antarCount int64
		database.DB.Model(&models.Trip{}).
			Joins("JOIN passengers ON passengers.trip_id = trips.id").
			Where("passengers.name = ? AND DATE(trips.trip_date) = DATE(?)", passenger, tripTime).
			Where("trips.trip_type = ?", "antar").
			Count(&antarCount)

		fare := OneWayFare
		tripType := "satu arah"
		if antarCount > 0 {
			fare = TwoWayFare
			tripType = "pulang pergi"
		}

		response += fmt.Sprintf("- %s (%s): Rp %d\n", passenger, tripType, fare)
	}

	return response
}
