package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Source string `env:"INGESTION_SOURCE" envDefault:"espn"`

	ESPNSiteBaseURL string `env:"ESPN_SITE_BASE_URL,required"`
	ESPNCoreBaseURL string `env:"ESPN_SITE_CORE_URL,required"`
	ESPNTimezone    string `env:"ESPN_TIMEZONE" envDefault:"America/New_York"`

	ESPNHTTPTimeout time.Duration `env:"ESPN_HTTP_TIMEOUT" envDefault:"10s"`

	// 	game poller
	ESPNPlaysInterval time.Duration `env:"ESPN_PLAYS_INTERVAL" envDefault:"2s"`

	// scheduler
	ESPNStatusInterval time.Duration `env:"ESPN_STATUS_INTERVAL" envDefault:"15s"`
	ESPNStartLeadTime  time.Duration `env:"ESPN_START_LEAD_TIME" envDefault:"1m"`

	ESPNDailyRefreshHour   int `env:"ESPN_DAILY_REFRESH_HOUR" envDefault:"4"`
	ESPNDailyRefreshMinute int `env:"ESPN_DAILY_REFRESH_MINUTE" envDefault:"0"`

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
