package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/robzlabz/angkot/internal/core/constants"
	"github.com/robzlabz/angkot/internal/core/ports"
	"github.com/robzlabz/angkot/pkg/logging"
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

func (s *botService) ProcessDeparture(driverName string, passengers []string) (string, error) {
	logging.Info(fmt.Sprintf("[Service][ProcessDeparture]Processing departure for %s with passengers: %v", driverName, passengers))
	// Ekstrak nama driver dari baris "Driver: [nama]"
	var driver string
	if strings.HasPrefix(driverName, "Driver:") {
		driver = strings.TrimSpace(strings.TrimPrefix(driverName, "Driver: "))
	}

	if driver == "" {
		return "", fmt.Errorf("nama driver tidak ditemukan")
	}

	// Filter array passengers untuk hanya mengambil nama santri
	var santriList []string
	for _, line := range passengers {
		if strings.HasPrefix(line, "-") {
			santri := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			santriList = append(santriList, santri)
		}
	}

	if len(santriList) == 0 {
		return "", fmt.Errorf("daftar santri tidak ditemukan")
	}

	if err := s.storage.SaveDeparture(driver, santriList); err != nil {
		log.Printf("[Service][ProcessDeparture]Error saving departure for %s: %v", driver, err)
		return "", fmt.Errorf("gagal menyimpan data keberangkatan: %v", err)
	}

	var response strings.Builder
	response.WriteString("✅ Keberangkatan berhasil dicatat\n\n")
	response.WriteString("Driver: " + driver + "\n")
	response.WriteString("Daftar Santri:\n")

	for _, passenger := range santriList {
		hasDeparture, _ := s.storage.HasDepartureToday(passenger)
		price := constants.SingleTripPrice
		if hasDeparture {
			price = constants.RoundTripPrice
		}
		response.WriteString(fmt.Sprintf("- %s (Rp %d)\n", passenger, price))
	}

	return response.String(), nil
}

func (s *botService) ProcessReturn(driverName string, passengers []string) (string, error) {
	// Ekstrak nama driver dari baris "Driver: [nama]"
	var driver string
	if strings.HasPrefix(driverName, "Driver:") {
		driver = strings.TrimSpace(strings.TrimPrefix(driverName, "Driver:"))
	}

	if driver == "" {
		return "", fmt.Errorf("nama driver tidak terdaftar")
	}

	// find driver name in database
	exists, err := s.storage.IsDriverExists(driver)
	if err != nil {
		return "", fmt.Errorf("gagal memeriksa keberadaan driver: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("driver tidak terdaftar")
	}

	// Filter array passengers untuk hanya mengambil nama santri
	var santriList []string
	for _, line := range passengers {
		if strings.HasPrefix(line, "-") {
			santri := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			santriList = append(santriList, santri)
		}
	}

	if len(santriList) == 0 {
		return "", fmt.Errorf("daftar santri tidak ditemukan")
	}

	if err := s.storage.SaveReturn(driver, santriList); err != nil {
		log.Printf("[Service][ProcessReturn]Error saving return for %s: %v", driver, err)
		return "", fmt.Errorf("gagal menyimpan data kepulangan: %v", err)
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

	return response.String(), nil
}

func (s *botService) GetTodayReport() (string, error) {
	var response strings.Builder
	response.WriteString("📊 Laporan Hari Ini\n")
	response.WriteString("================\n\n")

	// Ambil semua driver yang beroperasi hari ini
	drivers, err := s.storage.GetTodayDrivers()
	if err != nil {
		log.Printf("[Service][GetTodayReport]Error getting today's drivers: %v", err)
		return "", fmt.Errorf("gagal membuat laporan: %v", err)
	}

	totalIncome := 0
	for _, driver := range drivers {
		response.WriteString(fmt.Sprintf("🚗 Driver: %s\n", driver))

		// Ambil penumpang keberangkatan
		departurePassengers, err := s.storage.GetDeparturePassengers(driver)
		if err != nil {
			continue
		}

		// Ambil penumpang kepulangan
		returnPassengers, err := s.storage.GetReturnPassengers(driver)
		if err != nil {
			continue
		}

		// Hitung pendapatan driver
		driverIncome := 0
		response.WriteString("Penumpang:\n")

		// Proses penumpang keberangkatan
		for _, passenger := range departurePassengers {
			hasTwoTrips := false
			for _, returnPass := range returnPassengers {
				if returnPass == passenger {
					hasTwoTrips = true
					break
				}
			}

			price := constants.SingleTripPrice
			tripType := "Sekali Jalan"
			if hasTwoTrips {
				price = constants.RoundTripPrice
				tripType = "Pulang-Pergi"
			}

			response.WriteString(fmt.Sprintf("- %s (Rp %d - %s)\n", passenger, price, tripType))
			driverIncome += price
			totalIncome += price
		}

		// Tambahkan penumpang kepulangan yang belum tercatat
		for _, passenger := range returnPassengers {
			isNewPassenger := true
			for _, depPass := range departurePassengers {
				if depPass == passenger {
					isNewPassenger = false
					break
				}
			}

			if isNewPassenger {
				response.WriteString(fmt.Sprintf("- %s (Rp %d - Sekali Jalan)\n",
					passenger, constants.SingleTripPrice))
				driverIncome += constants.SingleTripPrice
				totalIncome += constants.SingleTripPrice
			}
		}

		response.WriteString(fmt.Sprintf("💰 Total Driver: Rp %d\n\n", driverIncome))
	}

	response.WriteString(fmt.Sprintf("💰 Total Pendapatan: Rp %d\n", totalIncome))
	return response.String(), nil
}
