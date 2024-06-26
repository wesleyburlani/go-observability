package config

import (
	pkg_config "github.com/wesleyburlani/go-observability/pkg/config"
)

type Config struct {
	ServiceName              string `mapstructure:"SERVICE_NAME" validate:"required"`
	ServiceVersion           string `mapstructure:"SERVICE_VERSION" validate:"required"`
	LogEnabled               bool   `mapstructure:"LOG_ENABLED"`
	LogLevel                 string `mapstructure:"LOG_LEVEL" validate:"required"`
	HttpAddress              string `mapstructure:"HTTP_ADDRESS" validate:"required"`
	GrpcAddress              string `mapstructure:"GRPC_ADDRESS" validate:"required"`
	KafkaHosts               string `mapstructure:"KAFKA_HOSTS" validate:"required"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtelExporterOtlpInsecure bool   `mapstructure:"OTEL_EXPORTER_OTLP_INSECURE"`
	DatabaseUrl              string `mapstructure:"DATABASE_URL" validate:"required"`
	AmqpUrl                  string `mapstructure:"AMQP_URL" validate:"required"`
}

func LoadDotEnvConfig(path string) (Config, error) {
	return pkg_config.LoadDotEnvConfig[Config](path)
}
