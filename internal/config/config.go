package config

import (
	pkg_config "github.com/wesleyburlani/go-rest/pkg/config"
)

type Config struct {
	ServiceName    string `mapstructure:"SERVICE_NAME" validate:"required"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION" validate:"required"`
	LogEnabled     bool   `mapstructure:"LOG_ENABLED"`
	LogLevel       string `mapstructure:"LOG_LEVEL" validate:"required"`
	HttpAddress    string `mapstructure:"HTTP_ADDRESS" validate:"required"`
}

func LoadDotEnvConfig(path string) (Config, error) {
	return pkg_config.LoadDotEnvConfig[Config](path)
}
