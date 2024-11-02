package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogging() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "timestamp"

	// Membuat console encoder untuk output di terminal
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Menggabungkan output ke file dan console
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, err error, fields ...zap.Field) {
	Logger.Error(msg, append(fields, zap.Error(err))...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}
