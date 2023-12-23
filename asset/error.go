package asset

import (
	"fmt"
	"spider-go/logger"

	"github.com/spf13/viper"
)

var e *RootError = &RootError{}

func LoadErrorCode(path string, fileNames ...string) *RootError {

	configViper := viper.New()
	configViper.AddConfigPath(path)
	for _, fileName := range fileNames {
		configViper.SetConfigName(fileName)
	}

	if err := configViper.ReadInConfig(); err != nil {
		logger.L().Error(fmt.Sprintf("unble to read config file: %+v", err))
	}

	configViper.Unmarshal(e)

	return e
}

func E() *RootError {
	return e
}
