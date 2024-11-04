package ports

type Storage interface {
	SaveDriver(name string, chatID int64) error
	GetDrivers(chatID int64) ([]string, error)
	IsDriverExists(name string, chatID int64) (bool, error)
	SavePassenger(name string, chatID int64) error
	GetPassengers(chatID int64) ([]string, error)
	SaveDeparture(driverName string, passengers []string, chatID int64) error
	SaveReturn(driverName string, passengers []string, chatID int64) error
	GetPassengerTripPrice(passengerName string, chatID int64) (int, error)
	GetDeparturePassengers(driverName string, chatID int64) ([]string, error)
	HasDepartureToday(passengerName string, chatID int64) (bool, error)
	GetTodayDrivers(chatID int64) ([]string, error)
	GetReturnPassengers(driverName string, chatID int64) ([]string, error)
}
