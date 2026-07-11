package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort int `env:"LIVE_GAME_SERVICE_PORT" envDefault:"8082"`

	RedisAddr string `env:"REDIS_ADDR" envDefault:"localhost:6379"`

	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`

	DynamoDBEndpoint string `env:"DYNAMODB_ENDPOINT" envDefault:"http://localhost:8000"`
	AWSRegion        string `env:"AWS_REGION" envDefault:"us-west-2"`
	AWSAccessKeyID   string `env:"AWS_ACCESS_KEY_ID" envDefault:"dummy"`
	AWSSecretKey     string `env:"AWS_SECRET_ACCESS_KEY" envDefault:"dummy"`
}

func Load() Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
