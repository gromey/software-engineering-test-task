package config

import (
	"fmt"

	"cruder/pkg/logger"
)

type Config struct {
	LogLevel logger.LogLevel `env:"LOG_LEVEL"`
	APIKey   string          `env:"API_KEY,m"`

	Host            string `env:"POSTGRES_HOST,m"`
	Port            uint16 `env:"POSTGRES_PORT,m"`
	User            string `env:"POSTGRES_USER,m"`
	Password        string `env:"POSTGRES_PASSWORD,m"`
	Database        string `env:"POSTGRES_DB,m"`
	PostgresSSLMode string `env:"POSTGRES_SSL_MODE,m"`
}

func (c *Config) GetPostgresDNS() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.PostgresSSLMode)
}
