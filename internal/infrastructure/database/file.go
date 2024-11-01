package database

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Add explicit interface implementation check

type FileDB struct {
	basePath string
}

func NewFileDB() *FileDB {
	// Create database directory if not exists
	if err := os.MkdirAll("database", 0755); err != nil {
		panic(err)
	}

	return &FileDB{
		basePath: "database",
	}
}

func (db *FileDB) SavePassenger(name string) error {
	filename := fmt.Sprintf("%s/passengers.txt", db.basePath)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("%s - %s\n", timestamp, name)); err != nil {
		return err
	}

	return nil
}

func (db *FileDB) GetPassengers() ([]string, error) {
	filename := fmt.Sprintf("%s/passengers.txt", db.basePath)

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []string{}, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var passengers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passengers = append(passengers, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return passengers, nil
}

func (db *FileDB) SaveDriver(name string) error {
	filename := fmt.Sprintf("%s/drivers.txt", db.basePath)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("%s - %s\n", timestamp, name)); err != nil {
		return err
	}

	return nil
}

func (db *FileDB) GetDrivers() ([]string, error) {
	filename := fmt.Sprintf("%s/drivers.txt", db.basePath)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []string{}, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var drivers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		drivers = append(drivers, scanner.Text())
	}

	return drivers, scanner.Err()
}

func (db *FileDB) IsDriverExists(name string) (bool, error) {
	drivers, err := db.GetDrivers()
	if err != nil {
		return false, err
	}

	for _, driver := range drivers {
		if strings.Contains(driver, name) {
			return true, nil
		}
	}
	return false, nil
}

func (db *FileDB) SaveDeparture(driver string, passengers []string) error {
	filename := fmt.Sprintf("%s/departures.txt", db.basePath)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf("=== Keberangkatan %s ===\n", timestamp)
	content += fmt.Sprintf("Driver: %s\n", driver)
	content += "Penumpang:\n"
	for _, passenger := range passengers {
		content += fmt.Sprintf("- %s\n", passenger)
	}
	content += "===========================\n\n"

	_, err = f.WriteString(content)
	return err
}

func (db *FileDB) SaveReturn(driver string, passengers []string) error {
	filename := fmt.Sprintf("%s/returns.txt", db.basePath)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf("=== Kepulangan %s ===\n", timestamp)
	content += fmt.Sprintf("Driver: %s\n", driver)
	content += "Penumpang:\n"
	for _, passenger := range passengers {
		content += fmt.Sprintf("- %s\n", passenger)
	}
	content += "===========================\n\n"

	_, err = f.WriteString(content)
	return err
}
