package logger

import "go.uber.org/zap/zapcore"

type Console struct {
	Enable        bool                   `json:"enable" yaml:"enable"`
	JsonFormat    bool                   `json:"jsonformat" yaml:"jsonformat"`
	EncoderConfig *zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	Level         string                 `json:"level" yaml:"level"`
}

type File struct {
	Enable        bool                   `json:"enable" yaml:"enable"`
	JsonFormat    bool                   `json:"jsonformat" yaml:"jsonformat"`
	EncoderConfig *zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	Level         string                 `json:"level" yaml:"level"`
	Path          string                 `json:"path" yaml:"path"`
	MaxSize       int                    `json:"maxsize" yaml:"maxsize"`
	MaxAge        int                    `json:"maxage" yaml:"maxage"`
}

type Config struct {
	ServiceName string  `json:"service" yaml:"service"`
	Console     Console `json:"console" yaml:"console"`
	File        File    `json:"file" yaml:"file"`
}
