package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevel = "LOG_LEVEL"
)

type Config struct {
	LogLevel zapcore.Level
}

func NewConfig() *Config {
	v := viper.New()

	v.SetDefault(LogLevel, zapcore.DebugLevel)

	v.AutomaticEnv()

	config := &Config{
		LogLevel: zapcore.Level(v.GetInt(LogLevel)),
	}

	return config
}
