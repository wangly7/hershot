package espn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wangly7/hershot/services/ingestion-service/internal/domain"
)

type SourceConfig struct {
	PollInterval            time.Duration
	StartLeadTime           time.Duration
	ScheduleRefreshInterval time.Duration
}

type ESPNSource struct {
	client Client
	config SourceConfig
}

func NewSource(
	client Client,
	config SourceConfig,
) *ESPNSource {
	return &ESPNSource{
		client: client,
		config: config,
	}
}

func (s *ESPNSource) Run(
	ctx context.Context,
	output chan<- domain.RawPlay,
) error {
	if s.client == nil {
		return errors.New("ESPN source client is nil")
	}

	if output == nil {
		return errors.New("ESPN source output channel is nil")
	}

	manager := NewPollManager(
		s.client,
		s.config.PollInterval,
		s.config.StartLeadTime,
		output,
	)

	scheduler := NewScheduler(
		s.client,
		manager,
		s.config.ScheduleRefreshInterval,
	)

	// This guarantees that all pending timers and running Pollers are stopped,
	// even if Scheduler.Run does not perform cleanup itself.
	defer manager.StopAll()

	scheduler.Run(ctx)

	if err := ctx.Err(); err != nil &&
		!errors.Is(err, context.Canceled) {
		return fmt.Errorf("run ESPN source: %w", err)
	}
	return nil
}
