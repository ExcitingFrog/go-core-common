package log

import (
	"github.com/ExcitingFrog/go-core-common/provider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

type Log struct {
	provider.IProvider

	Logger *zap.Logger
	Config *Config
}

func NewLog(config *Config) *Log {
	if config == nil {
		config = NewConfig()
	}

	return &Log{
		Config: config,
	}
}

func (l *Log) Run() error {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(l.Config.LogLevel),
		OutputPaths:      []string{"error.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	l.Logger = logger
	globalLogger = logger

	return nil
}

func (l *Log) Close() error {
	return l.Logger.Sync()
}

func Logger() *zap.Logger {
	return globalLogger
}
