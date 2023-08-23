package grpc_gateway

import (
	"github.com/spf13/viper"
)

const (
	GatawayPort string = "GATEWAY_PORT"
)

type Config struct {
	GatawayPort int
}

func NewConfig() *Config {
	v := viper.New()

	v.SetDefault(GatawayPort, 6061)

	v.AutomaticEnv()

	config := &Config{
		GatawayPort: v.GetInt(GatawayPort),
	}

	return config
}
