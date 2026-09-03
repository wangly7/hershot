package espn

import (
	"context"
	"errors"
	"log"
	"time"
)

const defaultScheduleRefreshInterval = 30 * time.Minute

// scheduler periodically loads today's ESPN scoreboard and update
// the timers managed by PollManager
type Scheduler struct {
	client  Client
	manager *PollManager

	refreshInterval time.Duration

	now func() time.Time
}

func NewScheduler(
	client Client,
	manager *PollManager,
	refreshInterval time.Duration,
) *Scheduler {
	if refreshInterval <= 0 {
		refreshInterval = defaultScheduleRefreshInterval
	}

	return &Scheduler{
		client:          client,
		manager:         manager,
		refreshInterval: refreshInterval,
		now:             time.Now,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	defer s.manager.StopAll()

	s.refreshAndLog(ctx)

	ticker := time.NewTicker(s.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.refreshAndLog(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) refreshAndLog(ctx context.Context) {
	err := s.refresh(ctx)
	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return
	}

	log.Printf("refresh ESPN schedule: %v", err)
}

// refresh explicitly requests the scoreboard for today's date
func (s *Scheduler) refresh(ctx context.Context) error {
	today := s.now()

	scoreboard, err := s.client.GetScoreboard(ctx, today)
	if err != nil {
		return err
	}

	games, err := MapScoreboard(scoreboard)
	if err != nil {
		return err
	}

	s.manager.SyncSchedule(ctx, games)

	return nil
}
