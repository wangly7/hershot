package config

type Config struct {
	HTTPPort int `env:"LIVE_GAME_SERVICE_PORT" envDefault:"8082"`

	RedisAddr string `env:"REDIS_ADDR" envDefault:"localhost:6379"`

	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`

	DynamoDBEndpoint string `env:"DYNAMODB_ENDPOINT" envDefault:"http://localhost:8000"`
	AWSRegion        string `env:"AWS_REGION" envDefault:"us-west-2"`
	AWSAccessKeyID   string `env:"AWS_ACCESS_KEY_ID" envDefault:"dummy"`
	AWSSecretKey     string `env:"AWS_SECRET_ACCESS_KEY" envDefault:"dummy"`
}
