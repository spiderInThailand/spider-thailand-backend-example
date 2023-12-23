package config

import (
	"fmt"
	"spider-go/logger"

	"github.com/spf13/viper"
)

var c *Root = &Root{}

func LoadConfig(path string, fileNames ...string) *Root {
	configViper := viper.New()
	configViper.AddConfigPath(path)
	for _, fileName := range fileNames {
		configViper.SetConfigName(fileName)
	}

	if err := configViper.ReadInConfig(); err != nil {
		logger.L().Error(fmt.Sprintf("unble to read config file: %+v", err))
	}

	configViper.Unmarshal(c)

	return c
}

func C() *Root {
	return c
}
