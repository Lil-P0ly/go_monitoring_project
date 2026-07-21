package logger

import (
	"fmt"

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
	if ZapLog == nil {
		return
	}
	ZapLog.Info(message, fields...)
}

func Infof(format string, args ...any) {
	if ZapLog == nil {
		return
	}
	ZapLog.Info(fmt.Sprintf(format, args...))
}

func Debug(message string, fields ...zap.Field) {
	if ZapLog == nil {
		return
	}
	ZapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	if ZapLog == nil {
		return
	}
	ZapLog.Error(message, fields...)
}

func Errorf(format string, args ...any) {
	if ZapLog == nil {
		return
	}
	ZapLog.Error(fmt.Sprintf(format, args...))
}

func Fatal(message string, fields ...zap.Field) {
	if ZapLog == nil {
		return
	}
	ZapLog.Fatal(message, fields...)
}

func Fatalf(format string, args ...any) {
	if ZapLog == nil {
		return
	}
	ZapLog.Fatal(fmt.Sprintf(format, args...))
}

func Sync() {
	if ZapLog == nil {
		return
	}
	ZapLog.Sync()
}
