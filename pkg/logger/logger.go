package logger

import (
	"log"
	"sync"
)

var (
	logger Logger
	once   sync.Once
)

const (
	INFO  = "info"
	ERROR = "error"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func Init(config Config) Logger {
	once.Do(func() {
		l, err := newZapLogger(config)
		if err != nil {
			log.Panic("failed to initialize logger")
		}

		logger = l
	})

	return logger
}

func Info(msg string) {
	logger.Info(msg)
}

func Error(msg string) {
	logger.Error(msg)
}
