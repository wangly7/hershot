package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort int `env:"INGESTION_SERVICE_PORT" envDefault:"8083"`

	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`

	SimulationIntervalM int `env:"SIMULATION_INTERVAL_MS" envDefault:"1000"`
}

func Load() Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
