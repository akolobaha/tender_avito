package config

import "github.com/ilyakaznacheev/cleanenv"

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

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	return c, nil
}
