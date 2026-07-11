package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort int `env:"LEAGUE_SERVICE_PORT" envDefault:"8081"`

	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"hershot"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:"hershot"`
	PostgresDB       string `env:"POSTGRES_DB" envDefault:"hershot"`
	PostgresSSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

func Load() Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}

func (c Config) PostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@s:%s:%d/%s?sslmode=%s",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
}
