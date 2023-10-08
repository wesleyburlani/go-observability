package config

import pkg_config "github.com/wesleyburlani/go-rest/pkg/config"

type Config struct {
	HttpAddress string `mapstructure:"HTTP_ADDRESS" validate:"required"`
}

func LoadDotEnvConfig(path string) (Config, error) {
	return pkg_config.LoadDotEnvConfig[Config](path)
}
