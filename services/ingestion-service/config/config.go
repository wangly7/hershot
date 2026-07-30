package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Source string `env:"SOURCE" envDefault:"espn"`

	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`

	KafkaGameEventsTopic string `env:"KAFKA_GAME_EVENTS_TOPIC" envDefault:"game-events"`

	ESPNSiteBaseURL string        `env:"ESPN_SITE_BASE_URL" envDefault:"https://site.api.espn.com"`
	ESPNCoreBaseURL string        `env:"ESPN_CORE_BASE_URL" envDefault:"https://sports.core.api.espn.com"`
	ESPNPlaysLimits int           `env:"ESPN_PLAYS_LIMITS" envDefault:"1000"`
	ESPNHTTPTimeout time.Duration `env:"ESPN_HTTP_TIMEOUT" envDefault:"5s"`

	PollInterval            time.Duration `env:"PLAYS_INTERVAL" envDefault:"2s"`
	StartLeadTime           time.Duration `env:"START_LEAD_TIME" envDefault:"30s"`
	ScheduleRefreshInterval time.Duration `env:"SCHEDULE_REFRESH_INTERVAL" envDefault:"30m"`

	OutputBuffers int `env:"OUTPUT_BUFFERS" envDefault:"256"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse ingestion-service config: %w", err)
	}

	return cfg, nil
}
