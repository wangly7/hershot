package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wangly7/hershot/shared/events"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(
	brokers []string,
	topic string,
) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("Kafka brokers cannot be empty")
	}

	if topic == "" {
		return nil, errors.New("Kafka topic cannot be empty")
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		// wait until all in-sync replicas acknowledege the message.
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		return nil, fmt.Errorf("create Kafka producer: %w", err)
	}
	return &Producer{
		client: client,
		topic:  topic,
	}, nil
}

func (p *Producer) Publish(
	ctx context.Context,
	event events.GameEvent,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if event.StreamID == "" {
		return errors.New("game event stream ID cannot be empty")
	}

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf(
			"marshal game event eventId=%s: %w",
			event.EventID,
			err,
		)
	}

	record := &kgo.Record{
		Topic: p.topic,
		Key:   []byte(event.StreamID),
		Value: value,
		Headers: []kgo.RecordHeader{
			{
				Key:   "content-type",
				Value: []byte("application/json"),
			},
			{
				Key:   "event-type",
				Value: []byte("game.play.created"),
			},
			{
				Key:   "source",
				Value: []byte(event.Source),
			},
			{
				Key:   "stream-mode",
				Value: []byte(event.StreamMode),
			},
		},
	}
	result := p.client.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		return fmt.Errorf(
			"publish game event eventId=%s gameId=%s topic=%s: %w",
			event.EventID,
			event.GameID,
			p.topic,
			err,
		)
	}
	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}
