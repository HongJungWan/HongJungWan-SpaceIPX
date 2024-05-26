package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func newZapLogger(config Config) (Logger, error) {
	cores := []zapcore.Core{}

	if config.Console.Enable {
		cores = append(cores, getConsoleCore(config.Console))
	}

	if config.File.Enable {
		cores = append(cores, getFileCore(config.File))
	}

	if config.ServiceName == "" {
		config.ServiceName = filepath.Base(os.Args[0])
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := addPid(addService(addHostname(
		zap.New(
			combinedCore,
			zap.AddCallerSkip(2),
			zap.AddCaller()).Sugar(),
	), config.ServiceName))

	return &zapLogger{sugar: logger}, nil
}

func getConsoleCore(consoleCfg Console) zapcore.Core {
	return zapcore.NewCore(
		getEncoder(consoleCfg.JsonFormat, consoleCfg.EncoderConfig),
		zapcore.Lock(os.Stdout),
		getZapLevel(consoleCfg.Level))
}

func getFileCore(fileCfg File) zapcore.Core {
	return zapcore.NewCore(
		getEncoder(fileCfg.JsonFormat, fileCfg.EncoderConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename: fileCfg.Path,
			MaxSize:  fileCfg.MaxSize, // megabytes
			MaxAge:   fileCfg.MaxAge,
			Compress: true,
		}),
		getZapLevel(fileCfg.Level))
}

func addPid(logger *zap.SugaredLogger) *zap.SugaredLogger {
	return logger.With(zap.Int("pid", os.Getpid()))
}

func addService(logger *zap.SugaredLogger, name string) *zap.SugaredLogger {
	return logger.With(zap.String("service", name))
}

func addHostname(logger *zap.SugaredLogger) *zap.SugaredLogger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return logger.With(zap.String("hostname", hostname))
}

func getEncoder(isJSON bool, encoderConfig *zapcore.EncoderConfig) zapcore.Encoder {
	var newEncoderConfig zapcore.EncoderConfig
	if encoderConfig == nil {
		newEncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		newEncoderConfig = *encoderConfig
	}

	if encoderConfig == nil || encoderConfig.EncodeTime == nil {
		newEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	if encoderConfig == nil || encoderConfig.EncodeLevel == nil {
		newEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if encoderConfig == nil || encoderConfig.EncodeCaller == nil {
		newEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	if encoderConfig == nil || encoderConfig.TimeKey == "" {
		newEncoderConfig.TimeKey = "time"
	}

	if isJSON {
		return zapcore.NewJSONEncoder(newEncoderConfig)
	}

	return zapcore.NewConsoleEncoder(newEncoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case INFO:
		return zapcore.InfoLevel
	case ERROR:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *zapLogger) Info(msg string) {
	l.sugar.Info(msg)
}

func (l *zapLogger) Error(msg string) {
	l.sugar.Error(msg)
}
