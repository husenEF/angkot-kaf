package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/robzlabz/angkot/internal/core/constants"
	"github.com/robzlabz/angkot/internal/core/ports"
)

type botService struct {
	storage                 ports.Storage
	waitingForPassengerName map[int64]bool
	waitingForDriverName    map[int64]bool
}

func NewBotService(storage ports.Storage) ports.BotService {
	return &botService{
		storage:                 storage,
		waitingForPassengerName: make(map[int64]bool),
		waitingForDriverName:    make(map[int64]bool),
	}
}

func (s *botService) HandlePing() string {
	return "pong"
}

func (s *botService) HandlePassenger(chatID int64) string {
	s.waitingForPassengerName[chatID] = true
	return "Siapa yang akan menjadi penumpang?"
}

func (s *botService) AddPassenger(name string) error {
	if err := s.storage.SavePassenger(name); err != nil {
		log.Printf("[Service][AddPassenger]Error failed to save passenger: %v", err)
		return err
	}
	return nil
}

func (s *botService) IsWaitingForPassengerName(chatID int64) bool {
	return s.waitingForPassengerName[chatID]
}

func (s *botService) ClearWaitingStatus(chatID int64) {
	delete(s.waitingForPassengerName, chatID)
	delete(s.waitingForDriverName, chatID)
}

func (s *botService) GetPassengerList() (string, error) {
	passengers, err := s.storage.GetPassengers()
	if err != nil {
		log.Printf("[Service][GetPassengerList]Error failed to get passengers: %v", err)
		return "", err
	}

	if len(passengers) == 0 {
		return "Belum ada penumpang terdaftar", nil
	}

	response := "Daftar Penumpang:\n"
	response += strings.Join(passengers, "\n")
	return response, nil
}

func (s *botService) HandleDriver(chatID int64) string {
	s.waitingForDriverName[chatID] = true
	return "Siapa yang akan menjadi driver?"
}

func (s *botService) AddDriver(name string) error {
	if err := s.storage.SaveDriver(name); err != nil {
		log.Printf("[Service][AddDriver]Error failed to save driver: %v", err)
		return err
	}
	return nil
}

func (s *botService) GetDriverList() (string, error) {
	drivers, err := s.storage.GetDrivers()
	if err != nil {
		return "", err
	}

	if len(drivers) == 0 {
		return "Belum ada driver terdaftar", nil
	}

	response := "Daftar Driver:\n"
	response += strings.Join(drivers, "\n")
	return response, nil
}

func (s *botService) IsWaitingForDriverName(chatID int64) bool {
	return s.waitingForDriverName[chatID]
}

func (s *botService) ProcessDeparture(text string) (string, error) {
	lines := strings.Split(text, "\n")
	if len(lines) < 3 {
		return "Format tidak valid. Gunakan format:\nKeberangkatan\nDriver: [nama]\n- [penumpang1]\n- [penumpang2]", nil
	}

	if !strings.Contains(strings.ToLower(lines[0]), "keberangkatan") {
		return "Format tidak valid. Baris pertama harus 'Keberangkatan'", nil
	}

	driverLine := lines[1]
	if !strings.HasPrefix(driverLine, "Driver:") {
		return "Format tidak valid. Baris kedua harus dimulai dengan 'Driver:'", nil
	}

	driverName := strings.TrimSpace(strings.TrimPrefix(driverLine, "Driver:"))
	exists, err := s.storage.IsDriverExists(driverName)
	if err != nil {
		log.Printf("[Service][ProcessDeparture]Error checking driver existence: %v", err)
		return "", err
	}
	if !exists {
		return fmt.Sprintf("Driver %s tidak terdaftar dalam database", driverName), nil
	}

	var passengers []string
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "-") {
			return "Format tidak valid. Daftar penumpang harus dimulai dengan '-'", nil
		}
		passenger := strings.TrimSpace(strings.TrimPrefix(line, "-"))
		passengers = append(passengers, passenger)
	}

	if len(passengers) == 0 {
		return "Minimal harus ada satu penumpang", nil
	}

	err = s.storage.SaveDeparture(driverName, passengers)
	if err != nil {
		log.Printf("[Service][ProcessDeparture]Error saving departure: %v", err)
		return "", err
	}

	response := fmt.Sprintf("Keberangkatan berhasil dicatat!\nDriver: %s\nJumlah Penumpang: %d\n\nBiaya per penumpang:\n",
		driverName, len(passengers))

	for _, passenger := range passengers {
		tripCount, err := s.storage.GetPassengerTripPrice(passenger)
		if err != nil {
			return "", err
		}

		price := constants.SingleTripPrice
		if tripCount > 1 {
			price = constants.RoundTripPrice / 2
		}

		response += fmt.Sprintf("- %s: Rp %d\n", passenger, price)
	}

	return response, nil
}

func (s *botService) ProcessReturn(text string) (string, error) {
	lines := strings.Split(text, "\n")
	if len(lines) < 3 {
		return "Format tidak valid. Gunakan format:\nKepulangan\nDriver: [nama]\n- [penumpang1]\n- [penumpang2]", nil
	}

	if !strings.Contains(strings.ToLower(lines[0]), "kepulangan") {
		return "Format tidak valid. Baris pertama harus 'Kepulangan'", nil
	}

	driverLine := lines[1]
	if !strings.HasPrefix(driverLine, "Driver:") {
		return "Format tidak valid. Baris kedua harus dimulai dengan 'Driver:'", nil
	}

	driverName := strings.TrimSpace(strings.TrimPrefix(driverLine, "Driver:"))
	exists, err := s.storage.IsDriverExists(driverName)
	if err != nil {
		log.Printf("[Service][ProcessReturn]Error checking driver existence: %v", err)
		return "", err
	}
	if !exists {
		return fmt.Sprintf("Driver %s tidak terdaftar dalam database", driverName), nil
	}

	var passengers []string
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "-") {
			return "Format tidak valid. Daftar penumpang harus dimulai dengan '-'", nil
		}
		passenger := strings.TrimSpace(strings.TrimPrefix(line, "-"))
		passengers = append(passengers, passenger)
	}

	if len(passengers) == 0 {
		return "Minimal harus ada satu penumpang", nil
	}

	// Get departure passengers first
	departurePassengers, err := s.storage.GetDeparturePassengers(driverName)
	if err != nil {
		log.Printf("[Service][ProcessReturn]Error getting departure passengers: %v", err)
		return "", err
	}

	// Create map of return passengers for easy lookup
	returnPassengersMap := make(map[string]bool)
	for _, p := range passengers {
		returnPassengersMap[p] = true
	}

	// Find passengers who only departed
	var onlyDeparture []string
	for _, dp := range departurePassengers {
		if !returnPassengersMap[dp] {
			onlyDeparture = append(onlyDeparture, dp)
		}
	}

	err = s.storage.SaveReturn(driverName, passengers)
	if err != nil {
		log.Printf("[Service][ProcessReturn]Error saving return: %v", err)
		return "", err
	}

	response := fmt.Sprintf("Kepulangan berhasil dicatat!\nDriver: %s\nJumlah Penumpang: %d\n\nBiaya per penumpang:\n",
		driverName, len(passengers))

	totalAmount := 0
	for _, passenger := range passengers {
		tripCount, err := s.storage.GetPassengerTripPrice(passenger)
		if err != nil {
			return "", err
		}

		var price int
		var note string
		if tripCount > 1 {
			price = constants.RoundTripPrice
			note = "(PP)"
			totalAmount += price
		} else {
			price = constants.SingleTripPrice
			note = "(Sekali jalan)"
			totalAmount += price
		}

		response += fmt.Sprintf("- %s: Rp %d %s\n", passenger, price, note)
	}

	if len(onlyDeparture) > 0 {
		response += "\nPenumpang yang hanya berangkat:\n"
		for _, passenger := range onlyDeparture {
			response += fmt.Sprintf("- %s (hanya berangkat)\n", passenger)
		}
	}

	response += fmt.Sprintf("\nTotal pembayaran: Rp %d", totalAmount)
	return response, nil
}
