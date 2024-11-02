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

func (s *botService) ProcessDeparture(driverName string, passengers []string) string {
	if err := s.storage.SaveDeparture(driverName, passengers); err != nil {
		log.Printf("[Service][ProcessDeparture]Error saving departure for %s: %v", driverName, err)
	}
	// ... existing processing code ...

	// Hitung biaya per penumpang
	var response strings.Builder
	response.WriteString("✅ Keberangkatan berhasil dicatat\n\n")
	response.WriteString("Driver: " + driverName + "\n")
	response.WriteString("Daftar Santri:\n")

	for _, passenger := range passengers {
		hasDeparture, _ := s.storage.HasDepartureToday(passenger)
		price := constants.SingleTripPrice
		if hasDeparture {
			price = constants.RoundTripPrice
		}
		response.WriteString(fmt.Sprintf("- %s (Rp %d)\n", passenger, price))
	}

	return response.String()
}

func (s *botService) ProcessReturn(driverName string, passengers []string) string {
	// Ekstrak nama driver dari baris "Driver: [nama]"
	var driver string
	for _, line := range passengers {
		if strings.HasPrefix(line, "Driver:") {
			driver = strings.TrimSpace(strings.TrimPrefix(line, "Driver:"))
			break
		}
	}

	// Filter array passengers untuk hanya mengambil nama santri
	var santriList []string
	for _, line := range passengers {
		if strings.HasPrefix(line, "-") {
			santri := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			santriList = append(santriList, santri)
		}
	}

	if err := s.storage.SaveReturn(driver, santriList); err != nil {
		log.Printf("[Service][ProcessReturn]Error saving return for %s: %v", driver, err)
	}

	var response strings.Builder
	response.WriteString("✅ Kepulangan berhasil dicatat\n\n")
	response.WriteString("Driver: " + driver + "\n")
	response.WriteString("Daftar Santri:\n")

	for _, passenger := range santriList {
		hasDeparture, _ := s.storage.HasDepartureToday(passenger)
		price := constants.SingleTripPrice
		if hasDeparture {
			price = constants.RoundTripPrice - constants.SingleTripPrice
			response.WriteString(fmt.Sprintf("- %s (Rp %d - Pulang-Pergi)\n", passenger, constants.RoundTripPrice))
		} else {
			response.WriteString(fmt.Sprintf("- %s (Rp %d - Sekali jalan)\n", passenger, price))
		}
	}

	return response.String()
}
