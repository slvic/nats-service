package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DatabaseConfig struct {
	DBUser     string `envconfig:"POSTGRES_USER"`
	DBName     string `envconfig:"POSTGRES_NAME"`
	DBPassword string `envconfig:"POSTGRES_PASSWORD"`
	DBPort     string `envconfig:"POSTGRES_PORT"`
	SSLMode    string `envconfig:"POSTGRES_SSL_MODE"`
}

func NewDbConfig() (*DatabaseConfig, error) {
	var dbConfig DatabaseConfig
	err := envconfig.Process("postgres", &dbConfig)
	if err != nil {
		return nil, fmt.Errorf("could not process env file: %s", err.Error())
	}
	return &DatabaseConfig{
		DBUser:     dbConfig.DBUser,
		DBName:     dbConfig.DBName,
		DBPassword: dbConfig.DBPassword,
		DBPort:     dbConfig.DBPort,
		SSLMode:    dbConfig.SSLMode,
	}, nil
}
