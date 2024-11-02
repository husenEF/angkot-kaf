package ports

type BotService interface {
	HandlePing() string
	HandlePassenger(chatID int64) string
	AddPassenger(name string) error
	IsWaitingForPassengerName(chatID int64) bool
	ClearWaitingStatus(chatID int64)
	GetPassengerList() (string, error)
	HandleDriver(chatID int64) string
	AddDriver(name string) error
	GetDriverList() (string, error)
	IsWaitingForDriverName(chatID int64) bool
	ProcessDeparture(string, []string) (string, error)
	ProcessReturn(string, []string) (string, error)
	GetTodayReport() (string, error)
}
