package services

import (
	"fmt"
	"strings"

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
	return s.storage.SavePassenger(name)
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
	return s.storage.SaveDriver(name)
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
		return "", err
	}

	return fmt.Sprintf("Keberangkatan berhasil dicatat!\nDriver: %s\nJumlah Penumpang: %d", driverName, len(passengers)), nil
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

	err = s.storage.SaveReturn(driverName, passengers)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Kepulangan berhasil dicatat!\nDriver: %s\nJumlah Penumpang: %d", driverName, len(passengers)), nil
}
