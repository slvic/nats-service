package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DatabaseConfig struct {
	DBUser     string `envconfig:"NATS_SERVICE_DB_USER"`
	DBName     string `envconfig:"NATS_SERVICE_DB_NAME"`
	DBPassword string `envconfig:"NATS_SERVICE_DB_PASSWORD"`
	DBHost     string `envconfig:"NATS_SERVICE_DB_HOST"`
	DBPort     string `envconfig:"NATS_SERVICE_DB_PORT"`
	SSLMode    string `envconfig:"NATS_SERVICE_DB_SSL_MODE"`
}

func NewDbConfig() (*DatabaseConfig, error) {
	var dbConfig DatabaseConfig
	err := envconfig.Process("nats_service_db", &dbConfig)
	if err != nil {
		return nil, fmt.Errorf("could not process env file: %s", err.Error())
	}
	return &dbConfig, nil
}
