package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var Cfg Config

type Config struct {
	Database Database `config:"DATABASE"`
}

// Initialize search the passed path and try to load the config file
func Initialize(configPath string) {
	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigFile(configPath)
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logrus.Info("No config file found.", err)
	}

	err := v.UnmarshalExact(&Cfg, func(f *mapstructure.DecoderConfig) {
		f.TagName = "config"
	})

	if err != nil {
		logrus.Panicf("invalid configuration: %s", err)
	}
}
