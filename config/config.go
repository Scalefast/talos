package config

import (
	"github.com/spf13/viper"
)

func ParseConfig(config string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(config)
	return v, v.ReadInConfig()
}
