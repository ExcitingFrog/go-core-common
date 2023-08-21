package jaeger

import (
	"github.com/spf13/viper"
)

const (
	ServiceName string = "SERVICE_NAME"
	JaegerURI   string = "JAEGER_URI"
)

type Config struct {
	ServiceName string
	JaegerURI   string
}

func NewConfig() *Config {
	v := viper.New()

	v.AutomaticEnv()

	v.SetDefault(JaegerURI, "http://localhost:14268/api/traces")

	config := &Config{
		ServiceName: v.GetString(ServiceName),
		JaegerURI:   v.GetString(JaegerURI),
	}

	return config
}
