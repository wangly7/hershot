package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RedpandaBrokers []string `env:"REDPANDA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`
	GameEventsTopic string   `env:"GAME_EVENTS_TOPIC" envDefault:"game-events"`

	SimulationIntervalSeconds int `env:"SIMULATION_INTERVAL_SECONDS" envDefault:"5"`

	GameID     string `env:"SIMULATION_GAME_ID" envDefault:"01"`
	HomeTeamID string `env:"SIMULATION_HOME_TEAM_ID" envDefault:"123"`
	AwayTeamID string `env:"SIMULATION_AWAY_TEAM_ID" envDefault:"456"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse ingestion-service config: %w", err)
	}

	if len(cfg.RedpandaBrokers) == 0 {
		return Config{}, fmt.Errorf("at least one Redpanda broker is required")
	}

	if cfg.SimulationIntervalSeconds <= 0 {
		return Config{}, fmt.Errorf("simulation interval must be positive")
	}

	return cfg, nil
}
