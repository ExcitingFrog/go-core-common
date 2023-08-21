package pprof

import (
	"github.com/spf13/viper"
)

const (
	PortKey     = "PPROF_PORT"
	EndpointKey = "PPROF_ENDPOINT"
)

type Config struct {
	port     int
	endpoint string
}

func NewConfig() *Config {
	v := viper.New()

	v.SetDefault(PortKey, 8083)
	v.SetDefault(EndpointKey, "/debug/pprof")

	viper.AutomaticEnv()

	config := &Config{
		port:     v.GetInt(PortKey),
		endpoint: v.GetString(EndpointKey),
	}

	return config
}
