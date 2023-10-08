package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	HttpAddress string `mapstructure:"HTTP_ADDRESS" validate:"required"`
}

func LoadEnvConfig(path string) (Config, error) {
	config := Config{}
	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return Config{}, err
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		log.Fatalf("Missing required attributes %v\n", err)
	}

	return config, nil
}
