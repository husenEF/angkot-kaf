package ports

type Storage interface {
	SaveDriver(name string) error
	GetDrivers() ([]string, error)
	IsDriverExists(name string) (bool, error)
	SavePassenger(name string) error
	GetPassengers() ([]string, error)
	SaveDeparture(driverName string, passengers []string) error
	SaveReturn(driverName string, passengers []string) error
}
