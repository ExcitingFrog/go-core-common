package utrace

import (
	"github.com/spf13/viper"
)

const (
	ServiceName string = "SERVICE_NAME"
	UptraceDSN  string = "UPTRACE_DSN"
)

type Config struct {
	ServiceName string
	UptraceDSN  string
}

func NewConfig() *Config {
	v := viper.New()

	v.SetDefault(UptraceDSN, "http://project2_secret_token@localhost:14317/2")

	v.AutomaticEnv()

	config := &Config{
		ServiceName: v.GetString(ServiceName),
		UptraceDSN:  v.GetString(UptraceDSN),
	}

	return config
}
