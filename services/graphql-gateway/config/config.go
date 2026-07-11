package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort int `env:"GRAPHQL_GATEWAY_PORT" envDefault:"8080"`

	LeagueServiceURL string `env:"LEAGUE_SERVICE_URL" envDefault:"http://localhost:8081"`
	RedisAddr        string `env:"REDIS_ADDR" envDefault:"localhost:6379"`

	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`
}

func Load() Config {
	var cfg Config

	if err := env.Parse((&cfg)); err != nil {
		log.Fatalf("failed to load config %v", cfg)
	}

	return cfg
}
