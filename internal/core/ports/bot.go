package ports

type BotService interface {
	HandlePing() string
	HandlePassenger(chatID int64) string
	AddPassenger(name string, chatID int64) error
	IsWaitingForPassengerName(chatID int64) bool
	ClearWaitingStatus(chatID int64)
	GetPassengerList(chatID int64) (string, error)
	HandleDriver(chatID int64) string
	AddDriver(name string, chatID int64) error
	GetDriverList(chatID int64) (string, error)
	IsWaitingForDriverName(chatID int64) bool
	ProcessDeparture(driverName string, passengers []string, chatID int64) (string, error)
	ProcessReturn(driverName string, passengers []string, chatID int64) (string, error)
	GetTodayReport(chatID int64) (string, error)
	GetReportByDate(chatID int64, date string) (string, error)
}
