package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wangly7/hershot/services/ingestion-service/internal/event"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func New(
	brokers []string,
	topic string,
) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(topic),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &Producer{
		client: client,
		topic:  topic,
	}, nil
}

func (p *Producer) PublishGameEvent(
	ctx context.Context,
	gameEvent event.GameEvent,
) error {
	value, err := json.Marshal(gameEvent)
	if err != nil {
		return fmt.Errorf("marshal game event: %w", err)
	}

	record := &kgo.Record{
		Topic: p.topic,
		Key:   []byte(gameEvent.GameID),
		Value: value,
		Headers: []kgo.RecordHeader{
			{
				Key:   "eventType",
				Value: []byte(string(gameEvent.EventType)),
			},
			{
				Key:   "eventId",
				Value: []byte(gameEvent.EventID),
			},
		},
	}

	result := p.client.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		return fmt.Errorf("publish game event: %w", err)
	}

	return nil
}

func (p *Producer) Ping(ctx context.Context) error {
	if err := p.client.Ping(ctx); err != nil {
		return fmt.Errorf("ping Redpanda: %w", err)
	}
	return nil
}

func (p *Producer) Close(ctx context.Context) {
	p.client.Close()
}
