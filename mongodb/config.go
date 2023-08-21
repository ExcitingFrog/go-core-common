package mongodb

import (
	"time"

	"github.com/spf13/viper"
)

const (
	URI            = "MONGODB_URI"
	MaxPoolSizeKey = "MONGODB_MAX_POOL_SIZE"
	MaxIdleKey     = "MONGODB_MAX_IDLE_SECOND"
	TimeoutKey     = "MONGODB_TIMEOUT_SECOND"
)

type Config struct {
	URI         string
	MaxPoolSize uint64
	MaxIdle     time.Duration
	Timeout     time.Duration
}

func NewConfig() *Config {
	v := viper.New()

	v.AutomaticEnv()

	v.SetDefault(URI, "")
	v.SetDefault(MaxPoolSizeKey, 32)
	v.SetDefault(MaxIdleKey, 15)
	v.SetDefault(TimeoutKey, 10)

	config := &Config{
		URI:         v.GetString(URI),
		Timeout:     v.GetDuration(TimeoutKey) * time.Second,
		MaxPoolSize: v.GetUint64(MaxPoolSizeKey),
		MaxIdle:     v.GetDuration(MaxIdleKey) * time.Second,
	}

	return config
}
