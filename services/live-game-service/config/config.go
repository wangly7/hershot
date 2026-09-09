package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort int `env:"LIVE_GAME_SERVICE_PORT" envDefault:"8082"`

	RedisAddr string `env:"REDIS_ADDR" envDefault:"localhost:6379"`

	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`

	DynamoDBEndpoint        string `env:"DYNAMODB_ENDPOINT" envDefault:"http://localhost:8000"`
	DynamoDBGameEventsTable string `env:"DYNAMODB_GAME_EVENTS_TABLE" envDefault:"game_events"`
	AWSRegion               string `env:"AWS_REGION" envDefault:"us-west-2"`
}

func Load() Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
