package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type ConfigDatabase struct {
	User     string `envconfig:"NATS_SERVICE_DB_USER"`
	Name     string `envconfig:"NATS_SERVICE_DB_NAME"`
	Password string `envconfig:"NATS_SERVICE_DB_PASSWORD"`
	Host     string `envconfig:"NATS_SERVICE_DB_HOST"`
	Port     string `envconfig:"NATS_SERVICE_DB_PORT"`
	ModeSSL  string `envconfig:"NATS_SERVICE_DB_SSL_MODE"`
}

type ConfigNATS struct {
	Host string `envconfig:"NATS_SERVICE_NATS_HOST"`
	Port string `envconfig:"NATS_SERVICE_NATS_PORT"`
}

func NewConfigDB() (*ConfigDatabase, error) {
	var dbConfig ConfigDatabase
	err := envconfig.Process("nats_service_db", &dbConfig)
	if err != nil {
		return nil, fmt.Errorf("could not process database env: %s", err.Error())
	}
	return &dbConfig, nil
}

func NewConfigNATS() (*ConfigNATS, error) {
	var natsConfig ConfigNATS
	err := envconfig.Process("nats_service_nats", &natsConfig)
	if err != nil {
		return nil, fmt.Errorf("could not process nats server env: %s", err.Error())
	}
	return &natsConfig, nil
}
