package grpc

import (
	"github.com/spf13/viper"
)

const (
	ServerPort string = "SERVER_PORT"
)

type Config struct {
	ServerPort int
}

func NewConfig() *Config {
	v := viper.New()

	v.SetDefault(ServerPort, 6060)

	v.AutomaticEnv()

	config := &Config{
		ServerPort: v.GetInt(ServerPort),
	}

	return config
}
