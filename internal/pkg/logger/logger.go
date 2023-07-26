package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

type Log interface {
	Info(message string)
	Warning(message string)
	Error(message string, err error)
	Fatal(message string, err error)
	Panic(message string, err error)
}

func NewLogger(env string) *Logger {
	log, _ := zap.NewDevelopment(zap.AddCallerSkip(2))

	if env == "production" {
		log, _ = zap.NewProduction(zap.AddCallerSkip(2))
	}

	return &Logger{log}
}

func (l *Logger) Info(message string) {
	l.log.Info(message)
}

func (l *Logger) Warning(message string) {
	l.log.Warn(message)
}

func (l *Logger) Error(message string, err error) {
	l.log.Error(message)
}

func (l *Logger) Fatal(message string, err error) {
	l.log.Fatal(message)
}

func (l *Logger) Panic(message string, err error) {
	l.log.Panic(message, zap.Error(err))
}
