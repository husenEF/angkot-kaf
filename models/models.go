package models

import (
	"time"
	"gorm.io/gorm"
)

type Driver struct {
	gorm.Model
	Name     string
	Trips    []Trip
	Active   bool
}

type Trip struct {
	gorm.Model
	TripDate    time.Time `gorm:"index"` // Add index for better query performance
	DriverID    uint
	Driver      Driver
	TripType    string    // "antar", "jemput", or "full"
	Passengers  []Passenger
}

type Passenger struct {
	gorm.Model
	TripID      uint
	Name        string
	TripType    string    // "antar", "jemput", or "full"
	Amount      float64   // 10000 for one way, 15000 for return trip
}
