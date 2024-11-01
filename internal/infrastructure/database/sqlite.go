package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robzlabz/angkot/internal/core/ports"
)

// Add explicit interface implementation check
var _ ports.Storage = (*SQLiteDB)(nil)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB() (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", "database/angkot.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initializeTables(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func initializeTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS drivers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS passengers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS departures (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			driver_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(driver_id) REFERENCES drivers(id)
		)`,
		`CREATE TABLE IF NOT EXISTS departure_passengers (
			departure_id INTEGER,
			passenger_name TEXT,
			FOREIGN KEY(departure_id) REFERENCES departures(id)
		)`,
		`CREATE TABLE IF NOT EXISTS returns (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			driver_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(driver_id) REFERENCES drivers(id)
		)`,
		`CREATE TABLE IF NOT EXISTS return_passengers (
			return_id INTEGER,
			passenger_name TEXT,
			FOREIGN KEY(return_id) REFERENCES returns(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (db *SQLiteDB) SaveDriver(name string) error {
	query := "INSERT INTO drivers (name) VALUES (?)"
	_, err := db.db.Exec(query, name)
	return err
}

func (db *SQLiteDB) GetDrivers() ([]string, error) {
	query := "SELECT name, created_at FROM drivers ORDER BY created_at DESC"
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []string
	for rows.Next() {
		var name string
		var createdAt time.Time
		if err := rows.Scan(&name, &createdAt); err != nil {
			return nil, err
		}
		drivers = append(drivers, fmt.Sprintf("%s - %s", createdAt.Format("2006-01-02 15:04:05"), name))
	}
	return drivers, nil
}

func (db *SQLiteDB) IsDriverExists(name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM drivers WHERE name = ?)"
	err := db.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}

func (db *SQLiteDB) SavePassenger(name string) error {
	query := "INSERT INTO passengers (name) VALUES (?)"
	_, err := db.db.Exec(query, name)
	return err
}

func (db *SQLiteDB) GetPassengers() ([]string, error) {
	query := "SELECT name, created_at FROM passengers ORDER BY created_at DESC"
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []string
	for rows.Next() {
		var name string
		var createdAt time.Time
		if err := rows.Scan(&name, &createdAt); err != nil {
			return nil, err
		}
		passengers = append(passengers, fmt.Sprintf("%s - %s", createdAt.Format("2006-01-02 15:04:05"), name))
	}
	return passengers, nil
}

func (db *SQLiteDB) SaveDeparture(driverName string, passengers []string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var driverID int
	err = tx.QueryRow("SELECT id FROM drivers WHERE name = ?", driverName).Scan(&driverID)
	if err != nil {
		return err
	}

	// Check if departure exists today for this driver
	var existingDepartureID int64
	today := time.Now().Format("2006-01-02")
	err = tx.QueryRow(`
		SELECT id FROM departures 
		WHERE driver_id = ? AND date(created_at) = date(?)`,
		driverID, today).Scan(&existingDepartureID)

	if err == nil {
		// Delete existing passengers for today's departure
		_, err = tx.Exec("DELETE FROM departure_passengers WHERE departure_id = ?", existingDepartureID)
		if err != nil {
			return err
		}
		// Update existing departure timestamp
		_, err = tx.Exec("UPDATE departures SET created_at = CURRENT_TIMESTAMP WHERE id = ?", existingDepartureID)
		if err != nil {
			return err
		}

		// Add new passengers to existing departure
		for _, passenger := range passengers {
			_, err = tx.Exec("INSERT INTO departure_passengers (departure_id, passenger_name) VALUES (?, ?)",
				existingDepartureID, passenger)
			if err != nil {
				return err
			}
		}
	} else {
		// Create new departure
		result, err := tx.Exec("INSERT INTO departures (driver_id) VALUES (?)", driverID)
		if err != nil {
			return err
		}

		departureID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Add passengers for new departure
		for _, passenger := range passengers {
			_, err = tx.Exec("INSERT INTO departure_passengers (departure_id, passenger_name) VALUES (?, ?)",
				departureID, passenger)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (db *SQLiteDB) SaveReturn(driverName string, passengers []string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var driverID int
	err = tx.QueryRow("SELECT id FROM drivers WHERE name = ?", driverName).Scan(&driverID)
	if err != nil {
		return err
	}

	// Check if return exists today for this driver
	var existingReturnID int64
	today := time.Now().Format("2006-01-02")
	err = tx.QueryRow(`
		SELECT id FROM returns 
		WHERE driver_id = ? AND date(created_at) = date(?)`,
		driverID, today).Scan(&existingReturnID)

	if err == nil {
		// Delete existing passengers for today's return
		_, err = tx.Exec("DELETE FROM return_passengers WHERE return_id = ?", existingReturnID)
		if err != nil {
			return err
		}
		// Update existing return timestamp
		_, err = tx.Exec("UPDATE returns SET created_at = CURRENT_TIMESTAMP WHERE id = ?", existingReturnID)
		if err != nil {
			return err
		}

		// Add new passengers to existing return
		for _, passenger := range passengers {
			_, err = tx.Exec("INSERT INTO return_passengers (return_id, passenger_name) VALUES (?, ?)",
				existingReturnID, passenger)
			if err != nil {
				return err
			}
		}
	} else {
		// Create new return
		result, err := tx.Exec("INSERT INTO returns (driver_id) VALUES (?)", driverID)
		if err != nil {
			return err
		}

		returnID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Add passengers for new return
		for _, passenger := range passengers {
			_, err = tx.Exec("INSERT INTO return_passengers (return_id, passenger_name) VALUES (?, ?)",
				returnID, passenger)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
