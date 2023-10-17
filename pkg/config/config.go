package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func LoadDotEnvConfig[T interface{}](path string) (T, error) {
	var config T
	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		return config, err
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
