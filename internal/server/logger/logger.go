package logger

import (
	"go.uber.org/zap"
)

var ZapLog *zap.Logger

func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		lvl, _ = zap.ParseAtomicLevel("info")
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl // Устанавливаем уровень

	ZapLog, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

func Info(message string, fields ...zap.Field) {
	ZapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	ZapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	ZapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	ZapLog.Fatal(message, fields...)
}

func Sync() {
	ZapLog.Sync()
}
