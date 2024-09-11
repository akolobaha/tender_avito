package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	PostgresConn     string `env:"POSTGRES_CONNECTION"`
	PostgresJdbcUrl  string `env:"POSTGRES_JDBC_URL"`
	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`
}

var Cfg *Config

var ConnString string

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	Cfg = c
	return c, nil
}

func InitDbConnectionString(c *Config) {
	ConnString = fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
		c.PostgresUsername, c.PostgresPassword, c.PostgresDatabase, c.PostgresHost, c.PostgresPort,
	)
}
